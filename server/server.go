package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const CreateTable = `CREATE TABLE IF NOT EXISTS cotacoes (
	"id" INTEGER PRIMARY KEY AUTOINCREMENT,
	"data" TEXT,
	"valor" TEXT
);`

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

	apiURL := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	defer cancelAPI()

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
	message := "R$ " + cotacaoResponse.USDBRL.BID

	cotacaoBytes, err := json.Marshal(cotacaoResponse.USDBRL.BID)
	if err != nil {
		log.Printf("erro ao serializar os dados: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(message)

	ctxDB, cancelDB := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelDB()

	db, err := sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(CreateTable)
	if err != nil {
		log.Printf("erro ao criar a tabela cotacoes: %v", err)
		return
	}

	data := time.Now().Format("2006-01-02 15:04:05")

	insertSQL := `INSERT INTO cotacoes (data, valor) VALUES (?, ?)`
	_, err = db.ExecContext(ctxDB, insertSQL, data, message)
	if err != nil {
		log.Printf("erro ao executar a query: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(cotacaoBytes)
}
