package dto

type TicketOrderResponseSave struct {
	OrderID       int
	TicketDetails TicketResponse
	UserDetails   UserResponse
	Amount        int
	TotalPrice    float64
}

type TicketOrderResponse struct {
	OrderID       int
	TicketDetails TicketResponseTicketOrder
	UserDetails   UserResponseTicketOrder
	Amount        int
	TotalPrice    float64
}