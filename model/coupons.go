package model

import "gorm.io/gorm"

type Coupons struct {
	gorm.Model
	UserId int     `json:"userId"`
	Number float64 `json:"number"`
}
