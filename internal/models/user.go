package models

import "time"

type User struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Handle          string    `gorm:"unique;not null" json:"handle"`
	DisplayName     string    `json:"display_name"`
	AtCoderHandle string  `gorm:"column:atcoder_handle" json:"atcoder_handle"`
	
	SeasonPoints    float64   `json:"season_points"`
	CFPoints        float64   `json:"cf_points"`
	AtCoderPoints float64 `gorm:"column:atcoder_points" json:"atcoder_points"`
	Activity7D      int       `gorm:"column:activity7_d" json:"activity_7d"`
	
	CurrentRating   int       `json:"current_rating"`
	MaxRating       int       `json:"max_rating"`
	BaseSolvedCount int       `json:"base_solved_count"` 
	TotalSolved     int       `json:"total_solved"`      
	HiddenSolved    int       `json:"hidden_solved"`     
	StruggleCount   int       `json:"struggle_count"`
	
	ManualBonus     int       `json:"manual_bonus"` 
	
	RankTier        string    `json:"rank_tier"`         
	PeakWeeklyRating int      `json:"peak_weekly_rating"` 
	
	LastSyncedAt    time.Time `json:"last_synced_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
