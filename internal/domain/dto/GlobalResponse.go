package dto

import "time"

type GlobalResponse struct {
	StatusCode     int
	StatusDesc     string
	Message        string
	RequestCreated string
	ProcessTime time.Duration
	Data           interface{}
}

func NewGlobalResponse(globalResponse GlobalResponse) GlobalResponse {
	return GlobalResponse{
		StatusCode: globalResponse.StatusCode,
		StatusDesc: globalResponse.StatusDesc,
		Message:    globalResponse.Message,
		RequestCreated: globalResponse.RequestCreated,
		ProcessTime: globalResponse.ProcessTime,
		Data:       globalResponse.Data,
	}
}