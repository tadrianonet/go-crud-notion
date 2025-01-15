package dependencies

import (
	"go-crud-notion/internal/interfaces/handlers"
	interfaces "go-crud-notion/internal/interfaces/services"
	impl "go-crud-notion/internal/services"
	usecase "go-crud-notion/internal/usecases"
	"log"

	"go.uber.org/dig"
)

func Setup() *dig.Container {
	container := dig.New()

	if err := container.Provide(func() interfaces.UserService {
		return impl.NewUserServiceImpl()
	}); err != nil {
		log.Fatalf("Erro ao registrar UserService: %v", err)
	}

	if err := container.Provide(usecase.NewUserUseCase); err != nil {
		log.Fatalf("Erro ao registrar UserUseCase: %v", err)
	}

	if err := container.Provide(handlers.NewUserHandler); err != nil {
		log.Fatalf("Erro ao registrar UserHandler: %v", err)
	}

	return container
}
