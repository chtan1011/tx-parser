# tx-parser
Ethereum blockchain parser that will allow to query transactions for subscribed addresses
## How to run? 

    go run main.go parser.go

## API
/subscribe

**Request method** 
POST

**Request header**

    'Content-Type: application/json'
**Request body**

    {'address': 'transaction address'}


/transactions/{address}

**Request method** 
GET


## Example 
**Subscribe address** 

    curl -X POST http://localhost:8080/subscribe -H "Content-Type: application/json" -d '{"address": "0xABC123"}'
    
    output: 
    Successfully subscribed to address: 0xABC123

**Check transaction**

    curl http://localhost:8080/transactions/0xABC123 
    
    output: 
    [{"Hash":"sample_hash","From":"0x123","To":"0xABC123","Value":"1 ETH","BlockNumber":20756477}]
