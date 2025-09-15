package models

import (
	"time"
)

type Portfolio struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	ProjectURL  string    `json:"project_url"`
	Category    string    `json:"category"`
	Tags        string    `json:"tags"` // Comma separated tags
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}