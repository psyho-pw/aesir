package messages

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	ChannelId uint    `json:"channelId"`
	Content   string  `json:"content"`
	Timestamp float64 `json:"timestamp"`
}
