package users

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	SlackId string `json:"slackId"`
	Alias   string `json:"alias"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
}
