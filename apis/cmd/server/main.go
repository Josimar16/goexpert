package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/josimar16/goexpert/apis/configs"
	entity "github.com/josimar16/goexpert/apis/internal/entities"
	"github.com/josimar16/goexpert/apis/internal/infra/database"
	handlers "github.com/josimar16/goexpert/apis/internal/infra/http"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	configs, err := configs.LoadConfig("./")

	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&entity.User{}, &entity.Product{})
	productDB := database.NewProduct(db)
	userDB := database.NewUser(db)

	ProductHandler := handlers.NewProductHandler(*productDB)
	UserHandler := handlers.NewUserHandler(*userDB, configs.TokenAuth, configs.JWTExperesIn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/products", ProductHandler.CreateProduct)
	r.Get("/products", ProductHandler.GetProducts)
	r.Get("/products/{id}", ProductHandler.GetProduct)
	r.Put("/products/{id}", ProductHandler.UpdateProduct)
	r.Delete("/products/{id}", ProductHandler.DeleteProduct)

	r.Post("/users", UserHandler.CreateUser)
	r.Post("/sessions", UserHandler.AuthenticateUser)

	// Start the server
	http.ListenAndServe(":3333", r)
}
