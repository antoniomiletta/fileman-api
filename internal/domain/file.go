package domain

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model

	ID          string    `json:"id" gorm:""`
	ParentID    string    `json:"parentId"`
	Name        string    `json:"name"`
	Extension   string    `json:"extension"`
	Size        int64     `json:"size"`
	StoragePath string    `json:"storagePath"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
