package domain

type Event struct {
	EventID   int    `json:"EventID"`
	EventName string `json:"EventName" validate:"required"`
}