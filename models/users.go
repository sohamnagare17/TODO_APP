package models

type Users struct {
	Userid   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
