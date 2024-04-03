package clients

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	ClientName string `json:"clientName"`
}
