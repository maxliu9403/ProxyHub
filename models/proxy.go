package models

type Proxy struct {
	Meta
	IP         string `json:"IP" gorm:"column:ip;type:varchar(64);not null;index:uq_proxy,unique;comment:'IP地址'"`
	Port       int64  `json:"Port" gorm:"column:port;not null;index:uq_proxy,unique;comment:'端口'"`
	Username   string `json:"Username" gorm:"column:username;type:varchar(128);not null;comment:'用户名'"`
	Password   string `json:"Password" gorm:"column:password;type:varchar(128);not null;comment:'密码'"`
	ProxyType  string `json:"ProxyType" gorm:"column:proxy_type;type:varchar(128);not null;default:socks5;comment:'代理类型，比如socks5'"`
	GroupID    int64  `json:"GroupID" gorm:"column:group_id;not null;index;comment:'所属代理池组'"`
	Source     string `json:"Source" gorm:"column:source;type:varchar(64);not null;index;comment:'来源类型，例：pias5/711/ipfoxy'"` //  新增字段
	InUseCount int64  `json:"InUseCount" gorm:"column:inuse_count;not null;index;comment:'当前使用数'"`
}

type ProxyBrief struct {
	IP       string `json:"IP"`
	Port     int64  `json:"Port"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type ReleaseIPDetail struct {
	IP    string `json:"IP"`
	Count int    `json:"Count"`
}

type UnbindEmulator struct {
	BrowserID string `json:"BrowserID"`
	UUID      string `json:"UUID"`
}
type GroupReleaseResult struct {
	GroupName       string            `json:"GroupName"`
	MaxOnline       int               `json:"MaxOnline"`
	ReleaseIPDetail []ReleaseIPDetail `json:"ReleaseIPDetail"`
	UnbindEmulator  []UnbindEmulator  `json:"UnbindEmulator"`
}
