package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name      string `json:"name" gorm:"varchar(20);not null"`
	Telephone string `json:"telephone" gorm:"varchar(20);not null;unique"`
	Password  string `json:"password" gorm:"size:255;not null"`
}

type UserResult struct {
	ID         int
	Created_at string
	Updated_at string
	Deleted_at string
	Name       string
	Telephone  string
	Password   string
}
