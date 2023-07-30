package messages

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	ChannelId uint   `json:"channelId"`
	Timestamp string `json:"timestamp"`
}
