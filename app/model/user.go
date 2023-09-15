package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name      string `gorm:"varchar(20);not null"`
	Telephone string `gorm:"varchar(20);not null;unique"`
	Password  string `gorm:"size:255;not null"`
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
