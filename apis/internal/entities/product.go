package entity

import (
	entity "github.com/josimar16/goexpert/apis/pkg/entities"
)

type Product struct {
	ID        entity.ID `json:"id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	CreatedAt string    `json:"created_at"`
}

func Create(name string, price int) (*Product, error) {
	return &Product{
		ID:        entity.UniqueEntityID(),
		Name:      name,
		Price:     price,
		CreatedAt: string(date),
	}, nil
}

func (product *Product) Validate() error {
	if product.ID.String() == "" {
		return ErrIDIsRequired
	}
}
