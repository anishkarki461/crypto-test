package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Using gorilla mux for ease of use in query parameters
type CurrencyResp struct {
	ID          string `json:"id"`
	FullName    string `json:"fullName"`
	Ask         string `json:"ask"`
	Bid         string `json:"bid"`
	Last        string `json:"last"`
	Open        string `json:"open"`
	Low         string `json:"low"`
	High        string `json:"high"`
	FeeCurrency string `json:"feeCurrency"`
}

type Ticker struct {
	Symbol      string    `json:"symbol"`
	Ask         string    `json:"ask"`
	Bid         string    `json:"bid"`
	Last        string    `json:"last"`
	Open        string    `json:"open"`
	Low         string    `json:"low"`
	High        string    `json:"high"`
	Volume      string    `json:"volume"`
	VolumeQuote string    `json:"volumeQuote"`
	Timestamp   time.Time `json:"timestamp"`
}

type Symbol struct {
	ID                   string `json:"id"`
	BaseCurrency         string `json:"baseCurrency"`
	QuoteCurrency        string `json:"quoteCurrency"`
	QuantityIncrement    string `json:"quantityIncrement"`
	TickSize             string `json:"tickSize"`
	TakeLiquidityRate    string `json:"takeLiquidityRate"`
	ProvideLiquidityRate string `json:"provideLiquidityRate"`
	FeeCurrency          string `json:"feeCurrency"`
	MarginTrading        bool   `json:"marginTrading"`
	MaxInitialLeverage   string `json:"maxInitialLeverage"`
}

type Currency struct {
	ID                  string `json:"id"`
	FullName            string `json:"fullName"`
	Crypto              bool   `json:"crypto"`
	PayinEnabled        bool   `json:"payinEnabled"`
	PayinPaymentID      bool   `json:"payinPaymentId"`
	PayinConfirmations  int    `json:"payinConfirmations"`
	PayoutEnabled       bool   `json:"payoutEnabled"`
	PayoutIsPaymentID   bool   `json:"payoutIsPaymentId"`
	TransferEnabled     bool   `json:"transferEnabled"`
	Delisted            bool   `json:"delisted"`
	PayoutFee           string `json:"payoutFee"`
	PayoutMinimalAmount string `json:"payoutMinimalAmount"`
	PrecisionPayout     int    `json:"precisionPayout"`
	PrecisionTransfer   int    `json:"precisionTransfer"`
	LowProcessingTime   string `json:"lowProcessingTime"`
	HighProcessingTime  string `json:"highProcessingTime"`
	AvgProcessingTime   string `json:"avgProcessingTime"`
}

type AllCurrencies struct {
	Currencies []struct {
		ID          string `json:"id"`
		FullName    string `json:"fullName"`
		Ask         string `json:"ask"`
		Bid         string `json:"bid"`
		Last        string `json:"last"`
		Open        string `json:"open"`
		Low         string `json:"low"`
		High        string `json:"high"`
		FeeCurrency string `json:"feeCurrency"`
	} `json:"currencies"`
}

func main() {
	// Initialize mux router and register the symbol
	router := mux.NewRouter()
	router.HandleFunc("/currency/{symbol}", GetCurrency)
	// router.HandleFunc("/currency/all", GetAllCurrency)
	http.ListenAndServe(":8080", router)
}

func GetCurrency(w http.ResponseWriter, req *http.Request) {
	// Parse the query variable to fetch the symbol
	vars := mux.Vars(req)
	symbol := vars["symbol"]
	fmt.Println(symbol)
	// if symbol is empty, handle it

	currBody, err := getCurrency(symbol)
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	currencyBytes, _ := json.Marshal(currBody)

	w.Header().Set("Content-Type", "application/json")
	w.Write(currencyBytes)
}

// func GetAllCurrency(w http.ResponseWriter, req *http.Request) {
// 	// Parse the query variable to fetch the symbol

// 	api := `https://api.hitbtc.com/api/2/public/ticker`

// 	resp, err := http.Get(api)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// 	defer resp.Body.Close()

// 	decoder := json.NewDecoder(resp.Body)
// 	var tickerCypto []Ticker
// 	err = decoder.Decode(&tickerCypto)
// 	fmt.Println("symbol: ", tickerCypto[0].Symbol)

// 	currRespArray AllCurrencies
// 	for i := 0; i < len(tickerCypto); i++ {
// 		currResp, err := getCurrency(tickerCypto[i].Symbol)
//      currRespArray = append(currRespArray.Currencies , currResp)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 	}
// 	currencyBytes, _ := json.Marshal(currBody)

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(currencyBytes)
// }

func getCurrency(symbol string) (CurrencyResp, error) {

	api := ""
	if symbol == "" {
		api = `https://api.hitbtc.com/api/2/public/ticker`
	} else {
		api = fmt.Sprintf(`https://api.hitbtc.com/api/2/public/ticker/%v`, symbol)
	}
	resp, err := http.Get(api)
	if err != nil {
		return CurrencyResp{}, err

	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var tickerCypto Ticker
	err = decoder.Decode(&tickerCypto)
	fmt.Println("symbol: ", tickerCypto.Symbol)
	if err != nil {
		return CurrencyResp{}, err
	}

	// Get baseCurrency from the symbol
	if tickerCypto.Symbol == "" {
		api = `https://api.hitbtc.com/api/2/public/symbol`
	} else {
		api = fmt.Sprintf(`https://api.hitbtc.com/api/2/public/symbol/%v`, tickerCypto.Symbol)
	}

	resp, err = http.Get(api)
	if err != nil {
		return CurrencyResp{}, err
	}
	defer resp.Body.Close()

	decoder = json.NewDecoder(resp.Body)
	var symbolCypto Symbol
	err = decoder.Decode(&symbolCypto)
	fmt.Println("symbol: ", symbolCypto.BaseCurrency)
	if err != nil {
		return CurrencyResp{}, err
	}

	// Get FullName from Symbol
	if symbolCypto.BaseCurrency == "" {
		api = `https://api.hitbtc.com/api/2/public/currency`
	} else {
		api = fmt.Sprintf(`https://api.hitbtc.com/api/2/public/currency/%v`, symbolCypto.BaseCurrency)
	}

	resp, err = http.Get(api)
	if err != nil {
		return CurrencyResp{}, err
	}
	defer resp.Body.Close()

	decoder = json.NewDecoder(resp.Body)
	var currencyCypto Currency
	err = decoder.Decode(&currencyCypto)
	fmt.Println("fullName: ", currencyCypto.FullName)
	if err != nil {
		return CurrencyResp{}, err
	}
	currBody := CurrencyResp{}
	currBody.Ask = tickerCypto.Ask
	currBody.Bid = tickerCypto.Bid
	currBody.FeeCurrency = symbolCypto.FeeCurrency
	currBody.FullName = currencyCypto.FullName
	currBody.High = tickerCypto.High
	currBody.ID = currencyCypto.ID
	currBody.Last = tickerCypto.Last
	currBody.Low = tickerCypto.Low
	currBody.Open = tickerCypto.Open

	return currBody, err

}
