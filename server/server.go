package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type USDBRL struct {
	BID string `json:"bid"`
}

type CotacaoResponse struct {
	USDBRL USDBRL `json:"USDBRL"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", Cotacao)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return
	}
}

func Cotacao(w http.ResponseWriter, r *http.Request) {
	ctxAPI, cancelAPI := context.WithTimeout(context.Background(), 200*time.Millisecond)
	//ctxDB, cancelDB := context.WithTimeout(context.Background(), 10*time.Millisecond)
	apiURL := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	defer cancelAPI()
	//defer cancelDB()

	cl := http.Client{}
	req, err := http.NewRequestWithContext(ctxAPI, "GET", apiURL, nil)
	if err != nil {
		log.Printf("erro ao criar a requisição: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := cl.Do(req.WithContext(ctxAPI))
	if err != nil {
		log.Printf("erro ao executar a requisição: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("erro ao ler a resposta da requisicão: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cotacaoResponse CotacaoResponse

	if err := json.Unmarshal(body, &cotacaoResponse); err != nil {
		log.Printf("erro ao fazer unmarshal da rersposta: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(cotacaoResponse.USDBRL.BID)
	cotacaoBytes, err := json.Marshal(cotacaoResponse.USDBRL.BID)
	if err != nil {
		log.Printf("erro ao serializar os dados: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(cotacaoBytes)
}
