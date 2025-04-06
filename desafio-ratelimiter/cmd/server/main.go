package main

import (
	"github.com/josimar16/goexpert/desafio-ratelimiter/config"
	"github.com/josimar16/goexpert/desafio-ratelimiter/router"
)

func main() {
	config.Init()
	router.Init()
}
