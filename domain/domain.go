package domain

import "time"

type User struct {
	ID        uint `gorm:"primaryKey"`
	ChatID    int64
	Username  string `gorm:"unique"`
	Role      string `gorm:"default:'user'"` // Может быть "user" или "admin"
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Message struct {
	ID        uint `gorm:"primaryKey"`
	Content   string
	UserID    uint      `gorm:"index"`
	Timestamp time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
