package models

type GetProfileByJWTRequestModel struct {
	Token string `header:"Authorization"`
}