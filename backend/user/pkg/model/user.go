package model

type User struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Email string `gorm:"uniqueIndex" json:"email"`
	Role  string `json:"role"` // "patient" or "doctor"
}
