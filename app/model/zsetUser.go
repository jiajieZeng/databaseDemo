package model

import "github.com/jinzhu/gorm"

type ZsetUser struct {
	gorm.Model
	Name  string  `json:"name" gorm:"varchar(20);not null"`
	Score float64 `json:"score" gorm:"float;not null"`
}
