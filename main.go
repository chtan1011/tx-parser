package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strings"
)

var parser *EthParser

// SubscribeHandler handles the /subscribe endpoint to subscribe to an Ethereum address
func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var req struct {
        Address string `json:"address"`
    }

    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    subscribed := parser.Subscribe(req.Address)
    if !subscribed {
        http.Error(w, "Address already subscribed", http.StatusConflict)
        return
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Successfully subscribed to address: %s\n", req.Address)
}

// TransactionsHandler handles the /transactions/{address} endpoint
func TransactionsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    pathParts := strings.Split(r.URL.Path, "/")
    if len(pathParts) != 3 {
        http.Error(w, "Invalid request path", http.StatusBadRequest)
        return
    }

    address := pathParts[2]

    transactions := parser.GetTransactions(address)
    if transactions == nil {
        http.Error(w, "No transactions found for address", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(transactions)
}

func main() {
    // Initialize the parser
    parser = NewEthParser()

    // Start polling the blockchain (this runs in the background)
    go PollBlockChain(parser)

    // Define the handlers
    http.HandleFunc("/subscribe", SubscribeHandler)
    http.HandleFunc("/transactions/", TransactionsHandler) // Matches /transactions/{address}

    // Start the HTTP server
    log.Println("Server is running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
