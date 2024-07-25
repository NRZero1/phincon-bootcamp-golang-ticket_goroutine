package domain

type Ticket struct {
	TicketID int
	EventID  int
	Name     string
	Price    float64
	Stock    int
}