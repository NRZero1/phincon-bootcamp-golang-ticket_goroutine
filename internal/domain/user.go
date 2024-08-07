package domain

type User struct {
	UserID      int     `json:"UserID"`
	Email       string  `json:"Email" validate:"required,email"`
	Name        string  `json:"Name" validate:"required"`
	PhoneNumber string  `json:"PhoneNumber" validate:"required,e164"`
	Balance     float64 `json:"Balance" validate:"required,number,gt=0"`
}