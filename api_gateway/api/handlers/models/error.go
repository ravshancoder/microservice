package models

// Error ...
type ResponseError struct {
	Error interface{} `json:"error"`
}

// StandardErrorModel ...
type StandartErrorModel struct {
	Error error `json:"error"`
}

type ErrorMessage struct {
	Message string  `json:"message"`
}