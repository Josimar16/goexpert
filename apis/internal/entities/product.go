package entity

import (
	"time"

	err "github.com/josimar16/goexpert/apis/internal/errors"
	entity "github.com/josimar16/goexpert/apis/pkg/entities"
)

type Product struct {
	ID        entity.ID `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

func NewProduct(name string, price float64) (*Product, error) {
	product := &Product{
		ID:        entity.UniqueEntityID(),
		Name:      name,
		Price:     price,
		CreatedAt: time.Now(),
	}

	err := product.Validate()

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (product *Product) Validate() error {
	if product.ID.String() == "" {
		return err.ErrIDIsRequired
	}
	if _, e := entity.ParseID(product.ID.String()); e != nil {
		return err.ErrIDInvalid
	}
	if product.Name == "" {
		return err.ErrNameIsRequired
	}
	if product.Price == 0 {
		return err.ErrPriceIsRequired
	}
	if product.Price < 0 {
		return err.ErrPriceInvalid
	}
	return nil
}
