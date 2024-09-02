package model

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	RequestMessage  string `json:"request_message"`
	ResponseMessage string `json:"response_message"`
	Url             string `json:"url"`
	Type            int    `json:"type"`
	UserId          int    `json:"user_id"`
}
