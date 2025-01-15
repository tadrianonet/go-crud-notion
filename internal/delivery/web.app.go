package webapp

import (
	"go-crud-notion/internal/delivery/dependencies"
	"go-crud-notion/internal/interfaces/handlers"
	"log"
	"net/http"
)

func Start() {
	container := dependencies.Setup()

	err := container.Invoke(func(userHandler *handlers.UserHandler) {
		http.HandleFunc("/users", userHandler.CreateUser)
		http.HandleFunc("/users/get", userHandler.GetUserByID)
		http.HandleFunc("/users/getall", userHandler.GetAllUsers)
		http.HandleFunc("/users/update", userHandler.UpdateUser)
		http.HandleFunc("/users/delete", userHandler.DeleteUserByPageID)
	})

	if err != nil {
		log.Fatalf("Erro ao resolver dependências: %v", err)
	}

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
