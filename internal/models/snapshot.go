package models

import "time"

type DailySnapshot struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `json:"user_id"`
	Handle        string    `json:"handle"`
	SeasonPoints  float64   `json:"season_points"`
	TotalSolved   int       `json:"total_solved"`
	Activity7D    int       `gorm:"column:activity7_d" json:"activity_7d"`
	CurrentRating int       `json:"current_rating"`
	RankTier      string    `json:"rank_tier"`
	SnapshotDate  time.Time `gorm:"uniqueIndex:idx_user_date;default:CURRENT_TIMESTAMP" json:"snapshot_date"`
	CreatedAt     time.Time `json:"created_at"`
}
