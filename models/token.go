package models

import "time"

type Token struct {
	Meta
	Token       string     `gorm:"column:token;type:varchar(64);primaryKey;comment:'访问令牌'" json:"Token"`
	Description string     `gorm:"column:description;type:varchar(128);not null;default:'';comment:'描述或备注信息'" json:"Description"`
	ExpireAt    *time.Time `gorm:"column:expire_at;comment:'过期时间，为空则永不过期'" json:"ExpireAt"`
}
