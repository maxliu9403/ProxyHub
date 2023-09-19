package models

import (
	"gorm.io/gorm"
)

type Meta struct {
	ID         int64          `json:"Id" gorm:"column:id;primaryKey"`
	Creator    string         `json:"Creator" gorm:"column:creator"`
	Updater    string         `json:"Updater" gorm:"column:updater"`
	Deleter    string         `json:"Deleter" gorm:"column:deleter"`
	CreateTime int64          `json:"CreateTime" gorm:"column:create_time;autoCreateTime"`
	UpdateTime int64          `json:"UpdateTime" gorm:"column:update_time;autoUpdateTime"` // 更新时间
	DeleteTime gorm.DeletedAt `json:"-" gorm:"column:delete_time"`
}
