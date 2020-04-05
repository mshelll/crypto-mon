# currency
This code base requires golang version 1.14 or above because modules are resolved based on go modules

To run the server make sure you have internet connection to download deps from github

        go run currency.go


To test the server

        curl -X GET "http://localhost:8080/currency/all"
        curl -X GET "http://localhost:8080/currency/ETHBTC"
        curl -X GET "http://localhost:8080/currency/BTCUSD"

CONGIURE allowed symbols

    The file currency.config in the same dir can be used to configure which all symbols can be allowed
    Even not configured ETHBSD and BTCUSD is allowed. During the start of server. We will fetch all 
    valid public symbols by HTTP request and validate all configured symbols in currency.config against 
    public symbols



How I Resolved Dependecies

From the currency dir run 
       
        go mod init currency

Above command will download and install all dependecy packages if not already present
For this execise I have used single external package "github.com/gorilla/websocket"
Socket market data was transmitted over websocket. To listen to socket market data we need above package
Alternative method is to use http request with a configurable polling time interval 






