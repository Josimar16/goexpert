package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/josimar16/goexpert/apis/configs"
	_ "github.com/josimar16/goexpert/apis/docs"
	entity "github.com/josimar16/goexpert/apis/internal/entities"
	"github.com/josimar16/goexpert/apis/internal/infra/database"
	handlers "github.com/josimar16/goexpert/apis/internal/infra/http"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	Product API with a JWT authentication
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Josimar Junior
//	@contact.url	http://www.buildcode.com.br
//	@contact.email	suporte@buildcode.com.br

//	@license.name	BuildCode
//	@license.url	http://www.buildcode.com.br

//	@host						localhost:3333
//	@BasePath					/
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization

func main() {
	configs, err := configs.LoadConfig(".")
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
	UserHandler := handlers.NewUserHandler(*userDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", configs.TokenAuth))
	r.Use(middleware.WithValue("JWTExperesIn", configs.JwtExperesIn))

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", ProductHandler.CreateProduct)
		r.Get("/", ProductHandler.GetProducts)
		r.Get("/{id}", ProductHandler.GetProduct)
		r.Put("/{id}", ProductHandler.UpdateProduct)
		r.Delete("/{id}", ProductHandler.DeleteProduct)
	})

	r.Post("/users", UserHandler.CreateUser)
	r.Post("/sessions", UserHandler.AuthenticateUser)

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:3333/docs/doc.json")))

	// Start the server
	http.ListenAndServe(":3333", r)
}
