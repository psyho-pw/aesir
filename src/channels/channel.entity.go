package channels

import "gorm.io/gorm"

type Channel struct {
	gorm.Model
	SlackId    string `gorm:"uniqueIndex" json:"slackId"`
	Name       string `json:"name"`
	Creator    string `json:"creator"`
	IsPrivate  bool   `json:"isPrivate"`
	IsArchived bool   `json:"isArchived"`
}
