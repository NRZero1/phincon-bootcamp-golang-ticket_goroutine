package dto

type UserResponse struct {
	UserID           int
	Email            string
	Name             string
	PhoneNumber      string
	RemainingBalance float64
}

type UserResponseTicketOrder struct {
	UserID      int
	Email       string
	Name        string
	PhoneNumber string
}