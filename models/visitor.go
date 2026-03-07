package models

import "time"

// Visitor stores each individual page visit.
type Visitor struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	IPAddress string    `json:"ip_address" gorm:"size:45;not null;index"`
	UserAgent string    `json:"user_agent" gorm:"type:text"`
	Path      string    `json:"path" gorm:"size:255;not null;index"`
	Referer   string    `json:"referer" gorm:"size:512"`
	Country   string    `json:"country" gorm:"size:100"`
	City      string    `json:"city" gorm:"size:100"`
	VisitedAt time.Time `json:"visited_at" gorm:"not null;index;default:now()"`
	CreatedAt time.Time `json:"created_at"`
}

// VisitorStats holds aggregated visitor statistics.
type VisitorStats struct {
	Today          int64            `json:"today"`
	ThisWeek       int64            `json:"this_week"`
	ThisMonth      int64            `json:"this_month"`
	ThisYear       int64            `json:"this_year"`
	Total          int64            `json:"total"`
	UniqueToday    int64            `json:"unique_today"`
	UniqueWeek     int64            `json:"unique_this_week"`
	UniqueMonth    int64            `json:"unique_this_month"`
	UniqueYear     int64            `json:"unique_this_year"`
	UniqueTotal    int64            `json:"unique_total"`
	DailyChart     []ChartDataPoint `json:"daily_chart"`
	MonthlyChart   []ChartDataPoint `json:"monthly_chart"`
	TopPages       []PageVisitCount `json:"top_pages"`
	RecentVisitors []Visitor        `json:"recent_visitors"`
}

// ChartDataPoint represents a single data point for visitor charts.
type ChartDataPoint struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// PageVisitCount represents visit counts per page.
type PageVisitCount struct {
	Path  string `json:"path"`
	Count int64  `json:"count"`
}
