package common

import "time"

type SQLModel struct {
	ID        int        `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}
