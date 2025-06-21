package models

import (
	"time"
)

type Token struct {
	Meta
	GroupID     int64      `json:"GroupID" gorm:"column:group_id;not null;index;comment:'所属代理池组'"`
	Token       string     `gorm:"column:token;type:varchar(64);primaryKey;comment:'访问令牌'" json:"Token"`
	Description string     `gorm:"column:description;type:varchar(128);not null;default:'';comment:'描述或备注信息'" json:"Description"`
	ExpireAt    *time.Time `gorm:"column:expire_at;comment:'过期时间，为空则永不过期'" json:"ExpireAt"`
}
