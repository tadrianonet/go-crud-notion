package entities

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Telefone string `json:"telefone"`
	PageId   string `json:"pageid"`
}
