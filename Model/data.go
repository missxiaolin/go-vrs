package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Data struct {
	Uid         int64  `bson:"uid" json:"uid"`                   // 用户ID
	IP          string `bson:"ip" json:"ip"`                     // IP地址
	FromUrl     string `bson:"from_url" json:"from_url"`         // 请求来源URL
	ActionUrl   string `bson:"action_url" json:"action_url"`     // 当前请求URL
	Method      string `bson:"method" json:"method"`             // 请求动作，如 增删改查
	Module      string `bson:"module" json:"module"`             // 项目模块名称
	SessionId   string `bson:"session_id" json:"session_id"`     // 会话ID
	Token       string `bson:"token" json:"token"`               // 令牌
	Content     string `bson:"content" json:"content"`           // http包请求体内容（前端发给后端的，非请求处理器的）
	PageId      string `bson:"page_id" json:"page_id"`           // 前端页面ID
	UserAgent   string `bson:"user_agent" json:"user_agent"`     // 用户浏览器标识
	UserChannel string `bson:"user_channel" json:"user_channel"` // 用户操作渠道
	TriggerTime int64  `bson:"trigger_time" json:"trigger_time"` // 操作触发时间（时间戳）
	ReceiveTime int64  `bson:"receive_time" json:"receive_time"` // 接收时间（时间戳）
	CreateTime  int64  `bson:"create_time" json:"create_time"`   // 日志生成时间（时间戳）
}

type Document struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"` // 日志ID
	Data
}
