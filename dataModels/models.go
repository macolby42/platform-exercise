package dataModels

type User struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
