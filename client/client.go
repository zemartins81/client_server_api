package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const CreateTable = `CREATE TABLE IF NOT EXISTS cotacoes (
	"id" INTEGER PRIMARY KEY AUTOINCREMENT,
	"data" TEXT,
	"valor" TEXT
);`

func main() {

	apiURL := "http://localhost:8080/cotacao"

	cl := http.Client{}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Printf("erro ao montar a requisição: %v", err)
		return
	}

	resp, err := cl.Do(req)
	if err != nil {
		log.Printf("erro ao fazer a requisição: %v", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("erro ao ler o corpo da resposta: %v", err)
		return
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Printf("erro ao criar o arquivo: %v", err)
		return
	}
	defer file.Close()

	message := fmt.Sprintf("Dolar: %s", body)

	_, err = file.Write([]byte(message))
	if err != nil {
		log.Printf("erro ao escrever no arquivo: %v", err)
		return
	}

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

	if resp.StatusCode != http.StatusOK {
		log.Printf("erro na resposta da API: %v", resp.Status)
		return
	}
}
