package domain

import "time"

type Folder struct {
	ID        string
	ParentID  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

