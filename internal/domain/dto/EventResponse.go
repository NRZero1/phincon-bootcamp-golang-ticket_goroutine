package dto

type EventResponse struct {
	EventID   int
	EventName string
}

func NewEventResponse(eventResponse EventResponse) EventResponse {
	return EventResponse{
		EventID:   eventResponse.EventID,
		EventName: eventResponse.EventName,
	}
}