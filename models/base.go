package models

import (
	"github.com/maxliu9403/ProxyHub/internal/types"
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

type GetListParams struct {
	types.BasicQuery          // Limit, Offset, Keyword, Order 等
	IPs              []string `json:"IPs,omitempty"`      // 多个 IP 精准匹配
	Ports            []int    `json:"Ports,omitempty"`    // 多端口匹配（如需）
	GroupIDs         []int64  `json:"GroupIDs,omitempty"` // 多组 ID 过滤
}

type GetTokenListParams struct {
	types.BasicQuery         // Limit, Offset, Keyword, Order 等
	GroupIDs         []int64 `json:"GroupIDs,omitempty"` // 多组 ID 过滤
}

type GetEmulatorListParams struct {
	types.BasicQuery          // Limit, Offset, Keyword, Order 等
	GroupIDs         []int64  `json:"GroupIDs,omitempty"` // 多组 ID 过滤
	UUIDS            []string `json:"UUIDS,omitempty"`    // uuid
	BrowserIDs       []string `json:"BrowserIDs,omitempty"`
}

type GetGroupListParams struct {
	types.BasicQuery          // Limit, Offset, Keyword, Order 等
	Name             []string `json:"Name,omitempty"` // 分组名称过滤
}
