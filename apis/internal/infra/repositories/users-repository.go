package repositories

import entity "github.com/josimar16/goexpert/apis/internal/entities"

type UsersRepository interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}
