package models

type Emulator struct {
	Meta
	BrowserID string `json:"BrowserID" gorm:"column:browser_id;uniqueIndex;comment:窗口ID"`
	UUID      string `json:"UUID" gorm:"column:uuid;uniqueIndex;comment:模拟器uuid"`
	GroupID   int64  `json:"GroupID" gorm:"index;column:group_id;comment:'分组ID'"`
	IP        string `json:"IP" gorm:"index;column:ip;comment:'IP'"`
}

type EmulatorBrief struct {
	BrowserID     string `json:"BrowserID"`
	UUID          string `json:"UUID"`
	IP            string `json:"IP"`
	SubscribeLink string `json:"SubscribeLink"`
}
