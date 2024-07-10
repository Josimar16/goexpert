package repositories

import entity "github.com/josimar16/goexpert/apis/internal/entities"

type ProductsRepository interface {
	Create(product *entity.Product) error
	FindAll(page, limit int, sort string) ([]entity.Product, error)
	FindByID(id string) (*entity.Product, error)
	Save(product *entity.Product) error
	Remove(id string) error
}
