package domain

type TicketOrder struct {
	OrderID    int     `json:"OrderID" validate:"required,number,gt=0"`
	TicketID   int     `json:"TicketID" validate:"required,number,gt=0"`
	UserID     int     `json:"UserID" validate:"required,number,gt=0"`
	Amount     int     `json:"Amount" validate:"required,number,gt=0"`
	TotalPrice float64 `json:"TotalPrice" validate:"required,number,gt=0"`
}