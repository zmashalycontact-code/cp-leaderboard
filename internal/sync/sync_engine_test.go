package sync

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/zmashaly/cp-leaderboard/internal/models"
)

const batchSize = 50

// ---------------------------------------------------------------------------
// Mock: UserRepository
// ---------------------------------------------------------------------------

type mockUserRepo struct {
	users         map[string]*models.User
	updateCalls   atomic.Int64
	lastSyncCalls atomic.Int64
}

func newMockUserRepo(users []models.User) *mockUserRepo {
	m := &mockUserRepo{users: make(map[string]*models.User, len(users))}
	for i := range users {
		u := users[i]
		m.users[strings.ToLower(u.Handle)] = &u
	}
	return m
}

func (m *mockUserRepo) Create(_ context.Context, u *models.User) error {
	m.users[strings.ToLower(u.Handle)] = u
	return nil
}

func (m *mockUserRepo) FindByID(_ context.Context, id uint) (*models.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepo) FindByHandle(_ context.Context, handle string) (*models.User, error) {
	u, ok := m.users[strings.ToLower(handle)]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

func (m *mockUserRepo) ListAll(_ context.Context) ([]models.User, error) {
	out := make([]models.User, 0, len(m.users))
	for _, u := range m.users {
		out = append(out, *u)
	}
	return out, nil
}

func (m *mockUserRepo) Update(_ context.Context, u *models.User) error {
	m.updateCalls.Add(1)
	m.users[strings.ToLower(u.Handle)] = u
	return nil
}

func (m *mockUserRepo) UpdateLastSynced(_ context.Context, id uint, t time.Time) error {
	m.lastSyncCalls.Add(1)
	return nil
}

// ---------------------------------------------------------------------------
// Mock: SnapshotRepository
// ---------------------------------------------------------------------------

type mockSnapshotRepo struct {
	upsertCalls atomic.Int64
	snapshots   []models.DailySnapshot
	mu          sync.Mutex
}

func (m *mockSnapshotRepo) Upsert(_ context.Context, s *models.DailySnapshot) error {
	m.upsertCalls.Add(1)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.snapshots = append(m.snapshots, *s)
	return nil
}

func (m *mockSnapshotRepo) FindByUserAndDateRange(_ context.Context, _ uint, _, _ time.Time) ([]models.DailySnapshot, error) {
	return nil, nil
}

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

var testLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

func makeTestUsers(handles ...string) []models.User {
	users := make([]models.User, len(handles))
	for i, h := range handles {
		users[i] = models.User{ID: uint(i + 1), Handle: h, CurrentRating: 1000}
	}
	return users
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestEngine_Run_HappyPath(t *testing.T) {
	handles := []string{"tourist", "Petr"}
	dbUsers := makeTestUsers(handles...)

	userRepo := newMockUserRepo(dbUsers)
	snapRepo := &mockSnapshotRepo{}

	engine := New(userRepo, snapRepo, nil, testLogger)
	engine.Run(context.Background())

	if upserts := snapRepo.upsertCalls.Load(); int(upserts) != len(handles) {
		t.Errorf("snapshot upserts = %d, want %d", upserts, len(handles))
	}
}

func TestEngine_Run_EmptyUserList(t *testing.T) {
	userRepo := newMockUserRepo(nil)
	snapRepo := &mockSnapshotRepo{}

	engine := New(userRepo, snapRepo, nil, testLogger)
	engine.Run(context.Background())

}

func TestEngine_Run_ContextCancelled(t *testing.T) {
	handles := []string{"tourist"}
	dbUsers := makeTestUsers(handles...)
	userRepo := newMockUserRepo(dbUsers)
	snapRepo := &mockSnapshotRepo{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	engine := New(userRepo, snapRepo, nil, testLogger)
	engine.Run(ctx)

}

func TestMockUserRepo_FindByID(t *testing.T) {
	users := makeTestUsers("alice", "bob")
	repo := newMockUserRepo(users)

	u, err := repo.FindByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.ID != 1 {
		t.Errorf("got ID %d, want 1", u.ID)
	}

	_, err = repo.FindByID(context.Background(), 999)
	if err == nil {
		t.Error("expected error for non-existent ID, got nil")
	}
}
