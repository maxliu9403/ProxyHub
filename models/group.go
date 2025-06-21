package models

type Groups struct {
	Meta
	Name        string `json:"Name" gorm:"column:name;uniqueIndex;comment:组名"`
	MaxOnline   int    `json:"MaxOnline" gorm:"index;column:max_online;comment:'该分组内的IP最大同时在线模拟器数'"`
	Description string `json:"Description" gorm:"index;column:description;comment:'描述'"`
}
