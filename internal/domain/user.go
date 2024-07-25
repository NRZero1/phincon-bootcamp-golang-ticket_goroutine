package domain

type User struct {
	UserID      int    `json:"UserID" validate:"required,number,gt=0"`
	Email       string `json:"Email" validate:"required,email"`
	Name        string `json:"Name" validate:"required"`
	PhoneNumber string `json:"PhoneNumber" validate:"required,e164"`
}