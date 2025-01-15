package services

import "go-crud-notion/internal/entities"

type UserService interface {
	Save(user *entities.User) error
	FindByID(id string) (*entities.User, error)
	FindAll() ([]*entities.User, error)
	UpdateUser(user *entities.User) error
	DeleteUserByPageID(id string) error
}
