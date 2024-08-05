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

// CreateProduct godoc
//
// @Summary      Create product
// @Description  Create a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateProductDTO true "product request"
// @Success      201
// @Failure      500  {object}  Error
// @Router       /products [post]
// @Security     ApiKeyAuth
func (productHandler *ProductHanlder) CreateProduct(w http.ResponseWriter, r *http.Request) {
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

// GetProduct godoc
//
// @Summary      Get product
// @Description  Get a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id    path      string true "product id"
// @Success      200   {object}  entity.Product
// @Failure      404   {object}  Error
// @Failure      500   {object}  Error
// @Router       /products/{id} [get]
// @Security     ApiKeyAuth
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
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// ListProducts godoc
//
// @Summary      List products
// @Description  List a products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page  query string false "page number"
// @Param        limit query string false "limit"
// @Success      200   {array} entity.Product
// @Failure      404   {object}  Error
// @Failure      500   {object}  Error
// @Router       /products [get]
// @Security     ApiKeyAuth
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

// UpdateProduct godoc
//
// @Summary      Update product
// @Description  Update a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id      path string true "product id"
// @Param        request body entity.Product true "product request"
// @Success      200
// @Failure      404   {object}  Error
// @Failure      500   {object}  Error
// @Router       /products/{id} [put]
// @Security     ApiKeyAuth
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
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	err = productHandler.ProductDB.Save(&product)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteProduct godoc
//
// @Summary      Delete product
// @Description  Delete a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id      path string true "product id"
// @Success      200
// @Failure      404   {object}  Error
// @Failure      500   {object}  Error
// @Router       /products/{id} [delete]
// @Security     ApiKeyAuth
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
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	err = productHandler.ProductDB.Remove(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
