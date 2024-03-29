//generated by lazy
//author: seanlan

package model

type PushMessageReq struct {
	BaseReq
	AppKey  string `form:"app_key" json:"app_key" binding:"required"` // app_key
	SendTo  string `form:"send_to" json:"send_to" binding:"required"` // 发送给谁
	Message string `form:"message" json:"message" binding:"required"` // 消息内容
	Nonce   string `form:"nonce" json:"nonce" binding:"required"`     // 随机数
	Sign    string `form:"sign" json:"sign" binding:"required"`       // 参数签名
}

type PushMessageResp struct {
}
