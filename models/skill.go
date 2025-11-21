package models

import (
	"time"
)

type Skill struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Level     string    `json:"level"`
	Score     int       `json:"score"`
	Category  string    `json:"category"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}