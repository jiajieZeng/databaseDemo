package model

type Persons struct {
	ID         int64  `gorm:"uniqueIndex;primarykey;bigint;not null"`
	Home       string `gorm:"varchar(40);not null"`
	Background string `gorm:"varchar(40);not null"`
}

type Infos struct {
	ID       int64  `gorm:"uniqueIndex;primarykey;bigint;not null"`
	Business string `gorm:"varchar(40);not null"`
	Address  string `gorm:"varchar(40);not null"`
}

type Belongings struct {
	ID          int64  `gorm:"uniqueIndex;primarykey;bigint;not null"`
	Cars        string `gorm:"varchar(40);not null"`
	Pets        string `gorm:"varchar(40);not null"`
	ClothesSize string `gorm:"varchar(40);not null"`
}

type RequestPerson struct {
	ID          int64  `json:"id" gorm:"primarykey;int;not null"`
	Home        string `json:"home" gorm:"varchar(40);not null"`
	Background  string `json:"background" gorm:"varchar(40);not null"`
	Business    string `json:"business" gorm:"varchar(40);not null"`
	Address     string `json:"address" gorm:"varchar(40);not null"`
	Cars        string `json:"cars" gorm:"primarykey;int;not null"`
	Pets        string `json:"pets" gorm:"varchar(40);not null"`
	ClothesSize string `json:"clothes_size" gorm:"varchar(40);not null"`
}
