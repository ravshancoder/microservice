package models

// Error ...
type Error struct {
	Message string `json:"message"`
}

// StandardErrorModel ...
type StandartErrorModel struct {
	Error error `json:"error"`
}
