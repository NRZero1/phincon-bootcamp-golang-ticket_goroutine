package dto

type Login struct {
	Email    string `json:"Email" validate:"email"`
	Password string `json:"Password"`
}