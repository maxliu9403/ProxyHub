package models

type ProxyGroups struct {
	Meta
	Name        string `json:"Name" gorm:"column:name;index;comment:'组名'"`
	MaxOnline   string `json:"MaxOnline" gorm:"index;column:max_online;comment:'最大同时在线模拟器数'"`
	Description string `json:"Description" gorm:"index;column:description;comment:'描述'"`
}
