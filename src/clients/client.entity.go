package clients

import (
	"time"
)
import "gorm.io/plugin/soft_delete"

type Client struct {
	ID         uint `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
	IsDel      soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt; uniqueIndex:idx_client_unique" json:"isDel"`
	ClientName string                `gorm:"uniqueIndex:idx_client_unique" json:"clientName"`
}
