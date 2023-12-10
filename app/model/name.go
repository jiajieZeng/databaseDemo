package model

type Names struct {
	ID       int64  `gorm:"uniqueIndex;primarykey;int;not null"`
	Name     string `gorm:"varchar(40);not null"`
	LastName string `gorm:"varchar(40);not null"`
}
