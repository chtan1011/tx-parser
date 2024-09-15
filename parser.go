package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
	"strconv"
    "sync"
    "time"
)

// Transaction structure to represent a transaction
type Transaction struct {
    Hash        string
    From        string
    To          string
    Value       string
    BlockNumber int
}

// Parser interface definition
type Parser interface {
    GetCurrentBlock() int
    Subscribe(address string) bool
    GetTransactions(address string) []Transaction
}

type EthParser struct {
    currentBlock int
    subscribers  map[string][]Transaction
    mu           sync.Mutex
}

func NewEthParser() *EthParser {
    return &EthParser{
        currentBlock: 0,
        subscribers:  make(map[string][]Transaction),
    }
}

func (ep *EthParser) GetCurrentBlock() int {
    ep.mu.Lock()
    defer ep.mu.Unlock()
    return ep.currentBlock
}

func (ep *EthParser) Subscribe(address string) bool {
    ep.mu.Lock()
    defer ep.mu.Unlock()

    if _, exists := ep.subscribers[address]; exists {
        return false // Already subscribed
    }
    ep.subscribers[address] = []Transaction{}
    return true
}

func (ep *EthParser) GetTransactions(address string) []Transaction {
    ep.mu.Lock()
    defer ep.mu.Unlock()
    return ep.subscribers[address]
}

func (ep *EthParser) UpdateCurrentBlock(block int) {
    ep.mu.Lock()
    defer ep.mu.Unlock()
    ep.currentBlock = block
}


func (ep *EthParser) AddTransaction(address string, tx Transaction) {
    ep.mu.Lock()
    defer ep.mu.Unlock()

    if _, exists := ep.subscribers[address]; exists {
        ep.subscribers[address] = append(ep.subscribers[address], tx)
    }
}

// JSONRPCRequest defines the structure for sending JSON-RPC requests
type JSONRPCRequest struct {
    Jsonrpc string      `json:"jsonrpc"`
    Method  string      `json:"method"`
    Params  interface{} `json:"params"`
    ID      int         `json:"id"`
}

// SendRPCRequest sends a JSON-RPC request to the Ethereum node and returns the response
func SendRPCRequest(method string, params interface{}) (map[string]interface{}, error) {
    rpcReq := JSONRPCRequest{
        Jsonrpc: "2.0",
        Method:  method,
        Params:  params,
        ID:      83,
    }

    reqBody, _ := json.Marshal(rpcReq)

    resp, err := http.Post("https://ethereum-rpc.publicnode.com", "application/json", bytes.NewBuffer(reqBody))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    return result, nil
}

func GetBlockNumber(currentblock int) (int, error) {
    result, err := SendRPCRequest("eth_blockNumber", []interface{}{})
    if err != nil {
        return 0, err
    }

    blockNumberHex := result["result"].(string)

	// Convert hex string to int
    blockNumber, err := strconv.ParseInt(blockNumberHex, 0, 0) 
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return currentblock , err
    }

    return int(blockNumber), nil
}

func PollBlockChain(parser *EthParser) {
    ticker := time.NewTicker(15 * time.Second)

    for range ticker.C {
		previousBlock := parser.GetCurrentBlock()
        currentBlock, err := GetBlockNumber(previousBlock)
        if err != nil {
            fmt.Println("Error fetching block number:", err)
            continue
        }

        if currentBlock > previousBlock {
            fmt.Println("Processing new block:", currentBlock)

            parser.UpdateCurrentBlock(currentBlock)
            for address := range parser.subscribers {
                tx := Transaction{
                    Hash:        "sample_hash",
                    From:        "0x123",
                    To:          address,
                    Value:       "1 ETH",
                    BlockNumber: currentBlock,
                }

                parser.AddTransaction(address, tx)
            }
        }
    }
}
