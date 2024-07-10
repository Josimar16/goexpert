package entity

import (
	"testing"

	error "github.com/josimar16/goexpert/apis/internal/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	product, err := NewProduct("Notebook", 1000.00)

	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.NotEmpty(t, product.ID)
	assert.NotEmpty(t, product.Price)
	assert.Equal(t, "Notebook", product.Name)
	assert.Equal(t, 1000.00, product.Price)
}

func TestProductWhenNameIsRequired(t *testing.T) {
	product, err := NewProduct("", 10.00)

	assert.Nil(t, product)
	assert.Equal(t, error.ErrNameIsRequired, err)
}

func TestProductWhenPriceIsRequired(t *testing.T) {
	product, err := NewProduct("Notebook", 0)

	assert.Nil(t, product)
	assert.Equal(t, error.ErrPriceIsRequired, err)
}

func TestProductWhenPriceIsInvalid(t *testing.T) {
	product, err := NewProduct("Notebook", -10)

	assert.Nil(t, product)
	assert.Equal(t, error.ErrPriceInvalid, err)
}

func TestProductValidate(t *testing.T) {
	product, err := NewProduct("Notebook", 1000.00)

	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.Nil(t, product.Validate())
}
