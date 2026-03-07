package repository

import (
	"backend-portfolio/models"
	"time"

	"gorm.io/gorm"
)

type visitorRepo struct{ db *gorm.DB }

// NewVisitorRepository returns a GORM-backed VisitorRepository.
func NewVisitorRepository(db *gorm.DB) VisitorRepository {
	return &visitorRepo{db: db}
}

// Record stores a new visitor entry.
func (r *visitorRepo) Record(visitor *models.Visitor) error {
	return r.db.Create(visitor).Error
}

// GetStats returns aggregated visitor statistics.
func (r *visitorRepo) GetStats() (*models.VisitorStats, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Start of week (Monday)
	weekday := now.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	startOfWeek := startOfDay.AddDate(0, 0, -int(weekday-time.Monday))

	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	stats := &models.VisitorStats{}

	// Total visits
	r.db.Model(&models.Visitor{}).Count(&stats.Total)
	r.db.Model(&models.Visitor{}).Where("visited_at >= ?", startOfDay).Count(&stats.Today)
	r.db.Model(&models.Visitor{}).Where("visited_at >= ?", startOfWeek).Count(&stats.ThisWeek)
	r.db.Model(&models.Visitor{}).Where("visited_at >= ?", startOfMonth).Count(&stats.ThisMonth)
	r.db.Model(&models.Visitor{}).Where("visited_at >= ?", startOfYear).Count(&stats.ThisYear)

	// Unique visitors (by IP)
	r.db.Model(&models.Visitor{}).Distinct("ip_address").Count(&stats.UniqueTotal)
	r.db.Model(&models.Visitor{}).Distinct("ip_address").Where("visited_at >= ?", startOfDay).Count(&stats.UniqueToday)
	r.db.Model(&models.Visitor{}).Distinct("ip_address").Where("visited_at >= ?", startOfWeek).Count(&stats.UniqueWeek)
	r.db.Model(&models.Visitor{}).Distinct("ip_address").Where("visited_at >= ?", startOfMonth).Count(&stats.UniqueMonth)
	r.db.Model(&models.Visitor{}).Distinct("ip_address").Where("visited_at >= ?", startOfYear).Count(&stats.UniqueYear)

	// Daily chart — last 30 days
	stats.DailyChart = r.getDailyChart(30)

	// Monthly chart — last 12 months
	stats.MonthlyChart = r.getMonthlyChart(12)

	// Top pages
	stats.TopPages = r.getTopPages(10)

	// Recent visitors (last 20)
	var recent []models.Visitor
	r.db.Order("visited_at DESC").Limit(20).Find(&recent)
	stats.RecentVisitors = recent

	return stats, nil
}

// GetDailyStats returns visitor counts grouped by day for the last N days.
func (r *visitorRepo) GetDailyStats(days int) ([]models.ChartDataPoint, error) {
	return r.getDailyChart(days), nil
}

// GetMonthlyStats returns visitor counts grouped by month for the last N months.
func (r *visitorRepo) GetMonthlyStats(months int) ([]models.ChartDataPoint, error) {
	return r.getMonthlyChart(months), nil
}

// GetTopPages returns the most visited pages.
func (r *visitorRepo) GetTopPages(limit int) ([]models.PageVisitCount, error) {
	return r.getTopPages(limit), nil
}

// GetRecentVisitors returns the most recent visitors.
func (r *visitorRepo) GetRecentVisitors(limit int) ([]models.Visitor, error) {
	var visitors []models.Visitor
	err := r.db.Order("visited_at DESC").Limit(limit).Find(&visitors).Error
	return visitors, err
}

// ── Private helpers ───────────────────────────────────────────

func (r *visitorRepo) getDailyChart(days int) []models.ChartDataPoint {
	since := time.Now().AddDate(0, 0, -days)
	var results []models.ChartDataPoint

	r.db.Model(&models.Visitor{}).
		Select("TO_CHAR(visited_at, 'YYYY-MM-DD') as date, COUNT(*) as count").
		Where("visited_at >= ?", since).
		Group("TO_CHAR(visited_at, 'YYYY-MM-DD')").
		Order("date ASC").
		Scan(&results)

	return results
}

func (r *visitorRepo) getMonthlyChart(months int) []models.ChartDataPoint {
	since := time.Now().AddDate(0, -months, 0)
	var results []models.ChartDataPoint

	r.db.Model(&models.Visitor{}).
		Select("TO_CHAR(visited_at, 'YYYY-MM') as date, COUNT(*) as count").
		Where("visited_at >= ?", since).
		Group("TO_CHAR(visited_at, 'YYYY-MM')").
		Order("date ASC").
		Scan(&results)

	return results
}

func (r *visitorRepo) getTopPages(limit int) []models.PageVisitCount {
	var results []models.PageVisitCount

	r.db.Model(&models.Visitor{}).
		Select("path, COUNT(*) as count").
		Group("path").
		Order("count DESC").
		Limit(limit).
		Scan(&results)

	return results
}
