package model

import (
	"gorm.io/gorm"
)

type RequestInfo struct {
	gorm.Model
	ApplicationID string `json:"application_id"`
	ChannelID     string `json:"channel_id"`
	SessionID     string `json:"session_id"`
	MJID          uint   `json:"mj_id"` // MjAccount的外键
	UserToken     string `json:"user_token"`
	GuildId       string `json:"guild_id"`
}

type MjAccount struct {
	gorm.Model
	Name     string `json:"name"`
	BotToken string `json:"bot_token"`
}
type MjInfo struct {
	gorm.Model
	Name        string        `json:"name"`
	BotToken    string        `json:"bot_token"`
	RequestInfo []RequestInfo `json:"requestInfo" gorm:"foreignKey:mj_id"`
}
