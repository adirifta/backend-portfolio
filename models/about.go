package models

import (
	"time"
)

type About struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Description1 string    `json:"description_singkat"`
	Description2 string    `json:"description_panjang"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	ImageURL    string    `json:"image_url"`
	LinkedinURL  string    `json:"linkedin_url"`
	GithubURL    string    `json:"github_url"`  
	InstagramURL   string    `json:"instagram_url"` 
	WebsiteURL   string    `json:"website_url"` 
	ResumeURL    string    `json:"resume_url"`  
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}