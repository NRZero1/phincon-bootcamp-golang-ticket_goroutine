package dto

type TicketResponse struct {
	TicketID     int `json:"TicketID" validate:"required,number,gt=0"`
	EventDetails EventResponse
	Name         string  `json:"Name" validate:"required"`
	Price        float64 `json:"Price" validate:"required,number,gt=0"`
	Stock        int     `json:"Stock" validate:"required,number,gt=0"`
	Type         string  `json:"Type" validate:"required"`
}