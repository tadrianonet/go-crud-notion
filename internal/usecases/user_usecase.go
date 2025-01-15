package usecases

import (
	"go-crud-notion/internal/entities"
	"go-crud-notion/internal/interfaces/services"
)

type UserUseCase struct {
	UserRepo services.UserService
}

func NewUserUseCase(repo services.UserService) *UserUseCase {
	return &UserUseCase{UserRepo: repo}
}

func (uc *UserUseCase) CreateUser(user *entities.User) error {
	return uc.UserRepo.Save(user)
}

func (uc *UserUseCase) GetUserByID(id string) (*entities.User, error) {
	return uc.UserRepo.FindByID(id)
}

func (uc *UserUseCase) GetAllUser() ([]*entities.User, error) {
	return uc.UserRepo.FindAll()
}

func (uc *UserUseCase) UpdateUser(user *entities.User) error {
	return uc.UserRepo.UpdateUser(user)
}

func (uc *UserUseCase) DeleteUserByPageID(id string) error {
	return uc.UserRepo.DeleteUserByPageID(id)
}
