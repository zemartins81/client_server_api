package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	ctxAPI, cancelAPI := context.WithTimeout(context.Background(), 300*time.Millisecond)

	defer cancelAPI()

	apiURL := "http://localhost:8080/cotacao"

	cl := http.Client{}

	req, err := http.NewRequestWithContext(ctxAPI, "GET", apiURL, nil)
	if err != nil {
		log.Printf("erro ao criar a requisição: %v\n", err)
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

	fmt.Println(message)

	_, err = file.Write([]byte(message))
	if err != nil {
		log.Printf("erro ao escrever no arquivo: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("erro na resposta da API: %v", resp.Status)
		return
	}
}
