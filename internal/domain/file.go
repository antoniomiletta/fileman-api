package domain

import "time"

type File struct {
	ID          string
	ParentID    string
	Name        string
	Extension   string
	Size        int64
	StoragePath string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
