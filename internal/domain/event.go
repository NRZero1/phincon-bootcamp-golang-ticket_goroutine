package domain

type Event struct {
	EventID   int    `json:"EventID" validate:"required,number"`
	EventName string `json:"EventName" validate:"required"`
}