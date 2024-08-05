package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/valyala/fastjson"
)

const brasilApiURL = "https://brasilapi.com.br/api/cep/v1/"
const viacepApiURL = "http://viacep.com.br/ws/"

type Address struct {
	Endpoint     string
	Cep          string
	State        string
	City         string
	Neighborhood string
	Street       string
}

func main() {
	cep01 := make(chan Address)
	cep02 := make(chan Address)

	go getAddressByBrasil(cep01, "01311-000")
	go getAddressByViaCep(cep02, "01311-000")

	select {
	case address := <-cep01:
		fmt.Printf("Endpoint: %s, CEP: %s, State: %s, City: %s", address.Endpoint, address.Cep, address.State, address.City)
	case address := <-cep02:
		fmt.Printf("Endpoint: %s, CEP: %s, State: %s, City: %s", address.Endpoint, address.Cep, address.State, address.City)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout")
	}
}

func getAddressByBrasil(c chan<- Address, cep string) {
	endpoint := brasilApiURL + cep

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)

	if err != nil {
		log.Fatalf("Internal Server Error: %v", err)
	}

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatalf("Internal Server Error: %v", err)
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

	var address Address

	address.Endpoint = string(endpoint)
	address.Cep = string(value.GetStringBytes("cep"))
	address.State = string(value.GetStringBytes("state"))
	address.City = string(value.GetStringBytes("city"))
	address.Neighborhood = string(value.GetStringBytes("neighborhood"))
	address.Street = string(value.GetStringBytes("street"))

	c <- address
}

func getAddressByViaCep(c chan<- Address, cep string) {
	endpoint := viacepApiURL + cep + "/json/"

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)

	if err != nil {
		log.Fatalf("Internal Server Error: %v", err)
	}

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatalf("Internal Server Error: %v", err)
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

	var address Address

	address.Endpoint = string(endpoint)
	address.Cep = string(value.GetStringBytes("cep"))
	address.State = string(value.GetStringBytes("uf"))
	address.City = string(value.GetStringBytes("localidade"))
	address.Neighborhood = string(value.GetStringBytes("bairro"))
	address.Street = string(value.GetStringBytes("logradouro"))

	c <- address
}
