package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/valyala/fastjson"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const awesomeapiURL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

type Cotacao struct {
	ID         int `gorm:"primaryKey"`
	Code       string
	Codein     string
	Name       string
	High       string
	Low        string
	VarBid     string
	PctChange  string
	Bid        string `json:"bid"`
	Ask        string
	Timestamp  string
	CreateDate string
	gorm.Model
}

type CotacaoReponse struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /cotacao", GetCotacaoHandler)

	log.Println("API running in port :8080")
	http.ListenAndServe(":8080", mux)
}

func GetCotacaoHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)

	defer cancel()

	cotacao, err := getCotacao(ctx, w)

	if err != nil {
		log.Fatalf("Internal Server Error: %v", err)
	}

	ctxDB, cancelDB := context.WithTimeout(context.Background(), 10*time.Millisecond)

	defer cancelDB()

	err = createCotacao(ctxDB, *cotacao)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Fatalf("Gateway Timeout: %v", err)
			w.WriteHeader(http.StatusGatewayTimeout)
		} else {
			log.Fatalf("Internal Server Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	data := map[string]string{
		"bid": cotacao.Bid,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

func getCotacao(ctx context.Context, w http.ResponseWriter) (*CotacaoReponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, awesomeapiURL, nil)

	if err != nil {
		log.Fatalf("Internal Server Error: %v", err)
	}

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Fatalf("Gateway Timeout: %v", err)
			w.WriteHeader(http.StatusGatewayTimeout)
		} else {
			log.Fatalf("Internal Server Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		log.Fatalf("Bad Request: %v", err)
	}

	defer response.Body.Close()

	var b fastjson.Parser

	value, err := b.Parse(string(body))

	if err != nil {
		log.Fatalf("Bad Request: %v", err)
	}

	data := value.Get("USDBRL").String()

	var cotacao CotacaoReponse

	if err := json.Unmarshal([]byte(data), &cotacao); err != nil {
		log.Fatalf("Bad Request: %v", err)
	}

	return &cotacao, nil
}

func createCotacao(ctx context.Context, cotacao CotacaoReponse) error {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	if err != nil {
		log.Fatalf("Erro ao abrir o banco de dados: %v", err)
		return err
	}

	db.AutoMigrate(&Cotacao{})

	return db.WithContext(ctx).Create(&Cotacao{
		Code:       cotacao.Code,
		Codein:     cotacao.Codein,
		Name:       cotacao.Name,
		High:       cotacao.High,
		Low:        cotacao.Low,
		VarBid:     cotacao.VarBid,
		PctChange:  cotacao.PctChange,
		Bid:        cotacao.Bid,
		Ask:        cotacao.Ask,
		Timestamp:  cotacao.Timestamp,
		CreateDate: cotacao.CreateDate,
	}).Error
}
