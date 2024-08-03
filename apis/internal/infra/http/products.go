package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/josimar16/goexpert/apis/internal/dto"
	entity "github.com/josimar16/goexpert/apis/internal/entities"
	"github.com/josimar16/goexpert/apis/internal/infra/database"
	pkg "github.com/josimar16/goexpert/apis/pkg/entities"
)

type ProductHanlder struct {
	ProductDB database.Product
}

func NewProductHandler(db database.Product) *ProductHanlder {
	return &ProductHanlder{ProductDB: db}
}

func (productHandler *ProductHanlder) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Create a product
	var body dto.CreateProductDTO
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := entity.NewProduct(body.Name, body.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = productHandler.ProductDB.Create(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (productHandler *ProductHanlder) GetProduct(w http.ResponseWriter, r *http.Request) {
	// Get a product
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := productHandler.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (productHandler *ProductHanlder) GetProducts(w http.ResponseWriter, r *http.Request) {
	// Get all products
	pageQuery := r.URL.Query().Get("page")
	limitQuery := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")

	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		page = 0
	}

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		limit = 10
	}

	products, err := productHandler.ProductDB.FindAll(page, limit, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (productHandler *ProductHanlder) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Update a product
	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product.ID, err = pkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = productHandler.ProductDB.FindByID(product.ID.String())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = productHandler.ProductDB.Save(&product)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (productHandler *ProductHanlder) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Delete a product
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := productHandler.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = productHandler.ProductDB.Remove(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
