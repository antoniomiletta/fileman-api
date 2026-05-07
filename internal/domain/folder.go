package domain

import "time"

type Folder struct {
	ID        string    `json:"id"`
	ParentID  string    `json:"parentId"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
