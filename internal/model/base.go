package model

type BaseReq struct {
	Platform  string `json:"platform" form:"platform" binding:"omitempty"`   // 平台
	Version   string `json:"version" form:"version" binding:"omitempty"`     // 版本
	TimeStamp int64  `json:"ts" form:"ts" binding:"omitempty"`               // 时间戳
	Channel   string `json:"channel" form:"channel" binding:"omitempty"`     // 渠道
	DeviceNo  string `json:"device_no" form:"device_no" binding:"omitempty"` // 设备号
	UserID    int64  // 用户ID
	ClientIP  string // 客户端IP
}
