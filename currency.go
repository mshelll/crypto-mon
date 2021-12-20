package main 

import (
	"net/http"
	"log"
	"strings"
	"github.com/gorilla/websocket"
	"encoding/json"
	"sync"
	"os"
	"bufio"
)

type Cache struct {

	data map[string][]byte
	mux sync.RWMutex
}

func (c *Cache) update(key string, val []byte) {

	c.mux.Lock()
	defer c.mux.Unlock()
	c.data[key] = val
}

func (c *Cache) fetch(key string) string {

	c.mux.RLock()
	defer c.mux.RUnlock()
	val, _ := c.data[key]
	return string(val)
}

var cache Cache
var valid_symbols map[string]struct{}

func currencyMonitor(symbol string) {

	url := "wss://api.hitbtc.com/api/2/ws"
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	defer c.Close()
	if err != nil {
		log.Fatal("webscoket dial:", err)
		return
	}

	message_tmpl := `{"method": "subscribeTicker", "params": {"symbol": "{symbol}" }, "id": 123}`
	request_msg := strings.NewReplacer("{symbol}", symbol).Replace(message_tmpl)

	if c.WriteMessage(websocket.TextMessage, []byte(request_msg)) != nil {

		log.Fatal("websocket write:", err)
		return
	}

	for {

		_, message, err := c.ReadMessage()
		if err != nil {
			log.Fatal("webscoket read:", err)
			return
		}

		var response map[string]interface{}
		err = json.Unmarshal(message, &response)
		if response["params"] != nil {
			params := response["params"]
			msg, _ := json.Marshal(params)
			cache.update(symbol, msg)
		}
	}
}

func currencyHandler(w http.ResponseWriter, req *http.Request) {

	url := req.URL.Path
	symbol := strings.TrimLeft(url, "/currency/")

	if symbol == "all" {

		response := make(map[string][]string)
		for symbol, _ := range(valid_symbols) {
			responseString := cache.fetch(symbol)
			response["currencies"] = append(response["currencies"], responseString)
		}
		json.NewEncoder(w).Encode(response)

	} else if _, found := valid_symbols[symbol]; found {

		response := cache.fetch(symbol)
		json.NewEncoder(w).Encode(response)
	} else {

		response := `{"code": 400, "message": "Invalid Symbol"`
		json.NewEncoder(w).Encode(response)
	}
}

func validateSymbols() {

	file, _ := os.Open("currency.config")
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var config_symbols []string

	for scanner.Scan() {
		config_symbols = append(config_symbols, scanner.Text())
	}

	type Symbol struct {

		Id string
	}

	var symbols []Symbol

	req_url := "https://api.hitbtc.com/api/2/public/symbol"  
	response, err := http.Get(req_url)
	defer response.Body.Close()
	if err != nil {
		log.Fatal("public symbol fetch :", err)
	}
	json.NewDecoder(response.Body).Decode(&symbols)

	pub_symbols := make(map[string]struct{})

	for _, symbol := range(symbols) {

		pub_symbols[symbol.Id] = struct{}{}
	}

	for _, symbol := range(config_symbols) {

		if _, found := pub_symbols[symbol]; found {

			valid_symbols[symbol] = struct{}{}
		} else {
			log.Println("configured symbol not valid :", symbol)
		}
	}
}

func init() {

	valid_symbols = map[string]struct{}{
		"ETHBTC": struct{}{},
		"BTCUSD": struct{}{},
	}

	cache.data = make(map[string][]byte)
}

func main() {

	validateSymbols()

	for symbol, _ := range(valid_symbols) {

		go currencyMonitor(symbol)
	}

	http.HandleFunc("/currency/", currencyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
} 