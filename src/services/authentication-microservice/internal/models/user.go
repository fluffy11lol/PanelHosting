package models

type User struct {
	ID       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	Active   bool   `json:"active" db:"active"`
	Token    string `json:"token"`
}
