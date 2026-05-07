package domain

type User struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"uniqueIndex;not null"`
	Password string `json:"-" gorm:"not null"`
}
