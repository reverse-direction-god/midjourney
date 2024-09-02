package model

import "gorm.io/gorm"

type UserInfo struct {
	gorm.Model
	TokenTime      string `json:"tokenTime"`
	Openid         string `json:"openid"`
	InvitationCode string `json:"invitationCode"`
}
