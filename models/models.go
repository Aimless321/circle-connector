package models

import "gorm.io/gorm"

type Error struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

type Body struct {
	ErrorField Error       `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func (b Body) Error() string {
	return b.ErrorField.Message
}

type MessageResponse struct {
	Message string `json:"message,omitempty"`
}

var Db *gorm.DB
