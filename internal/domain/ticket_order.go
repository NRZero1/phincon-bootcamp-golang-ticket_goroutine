package domain

type TicketOrder struct {
	OrderID    int
	TicketID   int
	UserID     int
	Amount     int
	TotalPrice float64
}