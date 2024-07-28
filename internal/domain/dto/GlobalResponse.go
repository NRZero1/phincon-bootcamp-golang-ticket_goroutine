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