package models

import (
	"time"
)

type Portfolio struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	Title       string          `json:"title" gorm:"not null"`
	Description string          `json:"description"`
	Category    string          `json:"category"`
	Tags        string          `json:"tags"`
	ProjectURL  string          `json:"project_url"`
	MediaFiles  []PortfolioMedia `json:"media_files" gorm:"foreignKey:PortfolioID;constraint:OnDelete:CASCADE;"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type PortfolioMedia struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	PortfolioID uint      `json:"portfolio_id"`
	Type        string    `json:"type"` // "image", "video"
	URL         string    `json:"url"`
	Thumbnail   string    `json:"thumbnail,omitempty"` // For video thumbnail
	OrderIndex  int       `json:"order_index" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
}