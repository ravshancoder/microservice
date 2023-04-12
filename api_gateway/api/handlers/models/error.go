package models

// Error ...
type ResponseError struct {
	Error interface{} `json:"error"`
}

type Error struct{
	Message string `json:"message"`
}

// StandardErrorModel ...
type StandartErrorModel struct {
	Error Error `json:"error"`
}

type ErrorMessage struct {
	Message string  `json:"message"`
}