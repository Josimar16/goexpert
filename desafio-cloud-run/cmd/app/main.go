package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/josimar16/goexpert/desafio-cloud-run/internal/infra/repo"
	"github.com/josimar16/goexpert/desafio-cloud-run/internal/infra/web"
	"github.com/josimar16/goexpert/desafio-cloud-run/internal/infra/web/webserver"
)

func ConfigureServer() *webserver.WebServer {
	webserver := webserver.NewWebServer(":8080")

	cepRepo := repo.NewCEPRepository()
	weatherRepo := repo.NewWeatherRepository(&http.Client{})

	open_weathermap_api_key := os.Getenv("OPEN_WEATHERMAP_API_KEY")
	if open_weathermap_api_key == "" {
		log.Fatal("Please, provide the OPEN_WEATHERMAP_API_KEY environment variable; Make sure you provide a valid api-key, otherwise it will not be possible to get and convert weather data")
	}

	webCEPHandler := web.NewWebCEPHandlerWithDeps(cepRepo, weatherRepo, os.Getenv("OPEN_WEATHERMAP_API_KEY"))
	webStatusHandler := web.NewWebStatusHandler()

	webserver.AddHandler("GET /cep/{cep}", webCEPHandler.Get)
	webserver.AddHandler("GET /status", webStatusHandler.Get)

	return webserver
}

func main() {
	webserver := ConfigureServer()
	fmt.Println("Starting web server on port", ":8080")
	webserver.Start()
}
