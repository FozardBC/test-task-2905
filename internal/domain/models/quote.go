package models

type Quote struct {
	Author string `json:"author" validate:"required,min=3,max=100,printascii"`
	Quote  string `json:"quote" validate:"required,min=3,max=500,printascii"`
}
