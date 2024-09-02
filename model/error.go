package model

import "gorm.io/gorm"

type Error struct {
	gorm.Model
	Message string `json:"message"`
}
