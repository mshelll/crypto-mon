# CRYPTO-MONITOR


# Dependencies 
Go > 1.12 <br>
(Gorilla Websockets)[https://github.com/gorilla/websocket]

# Build

```sh
 go build currency.go
```

# RUN

```sh
./currency
```

OR 

```sh
 go run currency.go
```

# TEST
```sh
        curl -X GET "http://localhost:8080/currency/all"
        curl -X GET "http://localhost:8080/currency/ETHBTC"
        curl -X GET "http://localhost:8080/currency/BTCUSD"
```

# CONFIGURE

The file currency.config in the same dir can be used to configure which all symbols can be allowed.By default ETHBSD and BTCUSD is monitored


## Resolve Dependecies

From the currency dir run 
```sh
go mod init currency
```






