package domain

type TicketOrder struct {
	OrderID    int `json:"OrderID"`
	TicketID   int `json:"TicketID" validate:"required,number,gt=0"`
	UserID     int `json:"UserID" validate:"required,number,gt=0"`
	Amount     int `json:"Amount" validate:"required,number,gt=0"`
	TotalPrice float64
}