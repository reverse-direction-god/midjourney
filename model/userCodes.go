package model

import "gorm.io/gorm"

type UserCodes struct {
	gorm.Model
	UserId     int `json:"userId"`
	UsedUserId int `json:"usedUserId"`
}
