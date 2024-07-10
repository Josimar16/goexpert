package database

import (
	"fmt"
	"math/rand"
	"testing"

	entity "github.com/josimar16/goexpert/apis/internal/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Product{})
	Product, _ := entity.NewProduct("Notebook", 1000.00)
	ProductDB := NewProduct(db)

	err = ProductDB.Create(Product)

	assert.Nil(t, err)

	var ProductFound entity.Product
	err = db.First(&ProductFound, "id = ?", Product.ID).Error

	assert.Nil(t, err)
	assert.Equal(t, Product.ID, ProductFound.ID)
	assert.Equal(t, Product.Name, ProductFound.Name)
	assert.Equal(t, Product.Price, ProductFound.Price)
}

func TestFindAllProducts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Product{})

	for i := 0; i < 24; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product %d", i+1), rand.Float64()*100)

		assert.NoError(t, err)
		db.Create(product)
	}

	productDB := NewProduct(db)

	products, err := productDB.FindAll(1, 10, "asc")

	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 10", products[9].Name)

	products, err = productDB.FindAll(2, 10, "asc")

	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 11", products[0].Name)
	assert.Equal(t, "Product 20", products[9].Name)

	products, err = productDB.FindAll(3, 10, "asc")

	assert.NoError(t, err)
	assert.Len(t, products, 4)
	assert.Equal(t, "Product 21", products[0].Name)
	assert.Equal(t, "Product 24", products[3].Name)
}

func TestFindProductByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Product{})
	Product, _ := entity.NewProduct("Notebook", 1000.00)
	ProductDB := NewProduct(db)

	err = ProductDB.Create(Product)

	assert.Nil(t, err)

	ProductFound, err := ProductDB.FindByID(Product.ID.String())

	assert.Nil(t, err)
	assert.Equal(t, Product.ID, ProductFound.ID)
	assert.Equal(t, Product.Name, ProductFound.Name)
	assert.Equal(t, Product.Price, ProductFound.Price)
}

func TestSaveProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Product{})
	product, _ := entity.NewProduct("Notebook", 1000.00)
	productDB := NewProduct(db)

	err = productDB.Create(product)

	assert.Nil(t, err)

	product.Name = "Macbook"
	err = productDB.Save(product)

	assert.Nil(t, err)

	product, err = productDB.FindByID(product.ID.String())

	assert.NoError(t, err)
	assert.Equal(t, "Macbook", product.Name)
}

func TestRemoveProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Product{})
	product, _ := entity.NewProduct("Notebook", 1000.00)
	productDB := NewProduct(db)

	err = productDB.Create(product)

	assert.Nil(t, err)

	err = productDB.Remove(product.ID.String())

	assert.NoError(t, err)

	_, err = productDB.FindByID(product.ID.String())

	assert.Error(t, err)
}
