package models

import (
	"time"
)

type QualificationType string

const (
	Education  QualificationType = "education"
	Experience QualificationType = "experience"
)

type Qualification struct {
	ID            uint               `json:"id" gorm:"primaryKey"`
	Type          QualificationType `json:"type"`
	Institution   string             `json:"institution"`
	Title         string             `json:"title"`
	Description   string             `json:"description"`
	StartDate     time.Time          `json:"start_date"`
	EndDate       *time.Time         `json:"end_date,omitempty"` // Pointer to allow null
	Current       bool               `json:"current"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}