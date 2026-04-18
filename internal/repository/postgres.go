package repository

import (
	"context"
	"time"

	"github.com/zmashaly/cp-leaderboard/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ---------------------------------------------------------------------------
// Interfaces (ports)
// ---------------------------------------------------------------------------

type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindByHandle(ctx context.Context, handle string) (*models.User, error)
	ListAll(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, u *models.User) error
	UpdateLastSynced(ctx context.Context, id uint, t time.Time) error
}

type SnapshotRepository interface {
	Upsert(ctx context.Context, s *models.DailySnapshot) error
	FindByUserAndDateRange(ctx context.Context, userID uint, from, to time.Time) ([]models.DailySnapshot, error)
}

// ---------------------------------------------------------------------------
// GORM implementations (adapters)
// ---------------------------------------------------------------------------

type postgresUserRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &postgresUserRepo{db: db}
}

func (r *postgresUserRepo) Create(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *postgresUserRepo) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var u models.User
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *postgresUserRepo) FindByHandle(ctx context.Context, handle string) (*models.User, error) {
	var u models.User
	if err := r.db.WithContext(ctx).
		Where("LOWER(handle) = LOWER(?)", handle).
		First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *postgresUserRepo) ListAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *postgresUserRepo) Update(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *postgresUserRepo) UpdateLastSynced(ctx context.Context, id uint, t time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("last_synced_at", t).Error
}

// ---------------------------------------------------------------------------

type postgresSnapshotRepo struct {
	db *gorm.DB
}

func NewSnapshotRepository(db *gorm.DB) SnapshotRepository {
	return &postgresSnapshotRepo{db: db}
}

func (r *postgresSnapshotRepo) Upsert(ctx context.Context, s *models.DailySnapshot) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "user_id"},
				{Name: "snapshot_date"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"season_points",
				"total_solved",
				"activity7_d"s,
				"current_rating",
				"rank_tier",
			}),
		}).
		Create(s).Error
}

func (r *postgresSnapshotRepo) FindByUserAndDateRange(
	ctx context.Context,
	userID uint,
	from, to time.Time,
) ([]models.DailySnapshot, error) {
	var snaps []models.DailySnapshot
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND snapshot_date BETWEEN ? AND ?", userID, from, to).
		Order("snapshot_date ASC").
		Find(&snaps).Error; err != nil {
		return nil, err
	}
	return snaps, nil
}
