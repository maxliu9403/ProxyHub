package models

type Proxy struct {
	Meta
	IP       string `json:"IP" gorm:"column:ip;type:varchar(64);not null;index:uq_proxy,unique;comment:'IP地址'"`
	Port     int64  `json:"Port" gorm:"column:port;not null;index:uq_proxy,unique;comment:'端口'"`
	Username string `json:"Username" gorm:"column:username;type:varchar(128);not null;comment:'用户名'"`
	Password string `json:"Password" gorm:"column:password;type:varchar(128);not null;comment:'密码'"`
	GroupID  int64  `json:"GroupID" gorm:"column:group_id;not null;index;comment:'所属代理池组'"`
	Enabled  bool   `json:"Enabled" gorm:"column:enabled;not null;default:true;index;comment:'是否启用'"`
	Source   string `json:"Source" gorm:"column:source;type:varchar(64);not null;index;comment:'来源类型，例：pias5/711/ipfoxy'"` //  新增字段
}
