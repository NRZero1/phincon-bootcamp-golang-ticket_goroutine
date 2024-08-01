package domain

type Ticket struct {
	TicketID int     `json:"TicketID"`
	EventID  int     `json:"EventID" validate:"required,number,gt=0"`
	Name     string  `json:"Name" validate:"required"`
	Price    float64 `json:"Price" validate:"required,number,gt=0"`
	Stock    int     `json:"Stock" validate:"required,number,gt=0"`
	Type     string  `json:"Type" validate:"required"`
}