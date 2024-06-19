package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const serverURL = "http://localhost:8080/cotacao"

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Second)

	defer cancel()

	cotacao, err := getCotacao(ctx)

	if err != nil {
		log.Fatalf("Internal Server Error: %v", err)
	}

	err = saveToFile(cotacao.Bid)

	if err != nil {
		log.Fatalf("Erro ao salvar o arquivo: %v", err)
	}

	fmt.Printf("O valor do dólar é R$ %s\n", cotacao.Bid)
}

func getCotacao(ctx context.Context) (*Cotacao, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL, nil)

	if err != nil {
		log.Fatalf("Internal Server Error: %v", err)
	}

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Fatalf("Gateway Timeout: %v", err)
		} else {
			log.Fatalf("Internal Server Error: %v", err)
		}
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		log.Fatalf("Bad Request: %v", err)
	}

	defer response.Body.Close()

	var cotacao Cotacao

	if err := json.Unmarshal(body, &cotacao); err != nil {
		log.Fatalf("Bad Request: %v", err)
	}

	return &cotacao, nil
}

func saveToFile(bid string) error {
	filePath := "cotacao.txt"
	_, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		file, err := os.Create(filePath)

		if err != nil {
			fmt.Println("Erro ao criar arquivo:", err)
			return err
		}

		defer file.Close()

		if _, err := file.WriteString(fmt.Sprintf("Dólar: %s\n", bid)); err != nil {
			fmt.Println("Erro ao escrever no arquivo:", err)
			return err
		}

	} else if err == nil {
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)

		if err != nil {
			fmt.Println("Erro ao abrir arquivo:", err)
			return err
		}

		defer file.Close()

		if _, err := file.WriteString(fmt.Sprintf("Dólar: %s\n", bid)); err != nil {
			fmt.Println("Erro ao escrever no arquivo:", err)
			return err
		}

	} else {
		fmt.Println("Erro ao verificar arquivo:", err)

		return err
	}
	return nil
}
