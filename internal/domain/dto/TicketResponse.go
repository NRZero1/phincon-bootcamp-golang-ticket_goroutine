package dto

type TicketResponse struct {
	TicketID     int
	EventDetails EventResponse
	Name         string
	Price        float64
	Stock        int
	Type         string
}

type TicketResponseTicketOrder struct {
	TicketID     int
	EventDetails EventResponse
	Name         string
	Price        float64
	Type         string
}