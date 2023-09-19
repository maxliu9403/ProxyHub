package models

type Demo struct {
	Meta
	User string `json:"User" gorm:"index;column:user"`
}
