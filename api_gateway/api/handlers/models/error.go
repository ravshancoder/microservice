package models

// Error ...
type ResponseError struct {
	Error interface{} `json:"error"`
}

// StandardErrorModel ...
type StandartErrorModel struct {
	Error error `json:"error"`
}
