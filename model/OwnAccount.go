package model

import "gorm.io/gorm"

type OwnAccount struct {
	gorm.Model
	BotToken       string `json:"botToken"`
	Application_id string `json:"applicationId"`
	Session_id     string `json:"sessionId"`
	Channel_id     string `json:"channelId"`
	User_token     string `json:"userToken"`
	Guild_id       string `json:"guildId"`
	User_id        uint   `json:"userId"`
}
type OwnAccountReq struct {
	Info           Queue  `json:"queue"`
	BotToken       string `json:"botToken"`
	Application_id string `json:"applicationId"`
	Session_id     string `json:"sessionId"`
	Channel_id     string `json:"channelId"`
	User_token     string `json:"userToken"`
	Guild_id       string `json:"guildId"`
}
