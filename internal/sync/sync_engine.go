package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	fhttp "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/redis/go-redis/v9"
	"github.com/zmashaly/cp-leaderboard/internal/models"
	"github.com/zmashaly/cp-leaderboard/internal/repository"
)

type CFSubmission struct {
	CreationTimeSeconds int64  `json:"creationTimeSeconds"`
	Verdict             string `json:"verdict"`
	Problem             struct {
		ContestId int    `json:"contestId"`
		Index     string `json:"index"`
		Rating    int    `json:"rating"`
	} `json:"problem"`
}

type SyncEngine struct {
	userRepo     repository.UserRepository
	snapshotRepo repository.SnapshotRepository
	rdb          *redis.Client
	logger       *slog.Logger
	startDate    int64
}

func New(ur repository.UserRepository, sr repository.SnapshotRepository, rdb *redis.Client, l *slog.Logger) *SyncEngine {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).Unix()
	return &SyncEngine{userRepo: ur, snapshotRepo: sr, rdb: rdb, logger: l, startDate: start}
}

func (s *SyncEngine) Run(ctx context.Context) {
	s.logger.Info("🔥 بدء دورة المزامنة...")

	users, err := s.userRepo.ListAll(ctx)
	if err != nil {
		s.logger.Error("❌ فشل جلب المستخدمين من DB", "err", err)
		return
	}

	if len(users) == 0 {
		s.logger.Warn("⚠️ لا يوجد مستخدمين في قاعدة البيانات")
		return
	}

	s.logger.Info("📋 بدء معالجة المستخدمين", "count", len(users))

	for i := range users {
		select {
		case <-ctx.Done():
			s.logger.Info("🛑 تم إيقاف المزامنة بسبب context cancellation")
			return
		default:
		}

		s.processUser(ctx, &users[i])
		time.Sleep(1 * time.Second)
	}

	s.logger.Info("✅ انتهت دورة المزامنة بنجاح", "users_synced", len(users))
}

func (s *SyncEngine) processUser(ctx context.Context, u *models.User) {

	scrapedTotal := s.scrapeTotalSolvedUltimate(u.Handle)
	if scrapedTotal > 0 {
		u.TotalSolved = scrapedTotal
	}

	apiCount, cfPts, hardAc, weekAct, curRating, apiHidden, peakRating, errCF := s.fetchStatusStats(u)

	if errCF == nil {
		if curRating > 0 {
			u.CurrentRating = curRating
		}
		u.PeakWeeklyRating = peakRating
		u.StruggleCount = hardAc

		totalSinceStart := u.TotalSolved - u.BaseSolvedCount
		if totalSinceStart < 0 {
			totalSinceStart = 0
		}

		publicCount := apiCount - apiHidden
		actualHidden := totalSinceStart - publicCount
		if actualHidden < 0 {
			actualHidden = 0
		}

		u.HiddenSolved = actualHidden

		missingFromApi := actualHidden - apiHidden
		if missingFromApi < 0 {
			missingFromApi = 0
		}

		if missingFromApi > 0 {
			cfPts += float64(missingFromApi * 4)
		}

	if u.ManualBonus != 0 {
				cfPts += u.ManualBonus
			}

		cfPts += float64(hardAc) * 0.5
		u.CFPoints = cfPts

		u.Activity7D = weekAct

		s.logger.Info("📊 تقرير حساب المتسابق",
			"handle", u.Handle,
			"weekAct_API", weekAct,
			"Missing_Sheets", missingFromApi,
			"Final_7D", u.Activity7D,
			"Total_Hidden", u.HiddenSolved,
		)

	} else {
		s.logger.Warn("⚠️ فشل CF API - تم الاحتفاظ بالبيانات القديمة", "handle", u.Handle, "err", errCF)
	}

	atcoderPts, errAC := s.fetchAtCoderStats(u.AtCoderHandle)
	if errAC == nil {
		u.AtCoderPoints = atcoderPts
	} else {
		s.logger.Warn("⚠️ فشل AtCoder API - تم الاحتفاظ بالبيانات القديمة", "handle", u.Handle, "err", errAC)
	}

	u.SeasonPoints = u.CFPoints + u.AtCoderPoints
	u.RankTier = s.getRankTier(u.SeasonPoints)
	u.LastSyncedAt = time.Now()

	s.userRepo.Update(ctx, u)

	snapshot := &models.DailySnapshot{
		UserID:        u.ID,
		Handle:        u.Handle,
		SeasonPoints:  u.SeasonPoints,
		TotalSolved:   u.TotalSolved,
		Activity7D:    u.Activity7D,
		CurrentRating: u.CurrentRating,
		RankTier:      u.RankTier,
		SnapshotDate:  time.Now().Truncate(24 * time.Hour),
		CreatedAt:     time.Now(),
	}
	s.snapshotRepo.Upsert(ctx, snapshot)

	if s.rdb != nil {
		userData, _ := json.Marshal(u)
		s.rdb.HSet(ctx, "leaderboard", u.Handle, userData)
	}
}

func (s *SyncEngine) scrapeTotalSolvedUltimate(handle string) int {
	urls := []string{
		fmt.Sprintf("https://codeforces.com/profile/%s?lang=en", handle),
		fmt.Sprintf("https://mirror.codeforces.com/profile/%s?lang=en", handle),
	}
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(15),
		tls_client.WithClientProfile(profiles.Chrome_120),
	}
	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return s.scrapeWithProxies(handle)
	}

	for _, u := range urls {
		req, _ := fhttp.NewRequest(fhttp.MethodGet, u, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0")
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if val := extractSolved(string(body)); val > 0 {
			return val
		}
	}
	return s.scrapeWithProxies(handle)
}

func (s *SyncEngine) scrapeWithProxies(handle string) int {
	encodedUrl := url.QueryEscape(fmt.Sprintf("https://codeforces.com/profile/%s?lang=en", handle))
	proxies := []string{
		"https://api.allorigins.win/raw?url=" + encodedUrl,
		"https://api.codetabs.com/v1/proxy?quest=" + encodedUrl,
	}
	client := &http.Client{Timeout: 15 * time.Second}
	for _, proxyUrl := range proxies {
		req, _ := http.NewRequest("GET", proxyUrl, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0")
		if resp, err := client.Do(req); err == nil && resp.StatusCode == 200 {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			if val := extractSolved(string(body)); val > 0 {
				return val
			}
		}
	}
	return 0
}

func extractSolved(html string) int {
	patterns := []string{
		`(?i)counterValue">(\d+)</div>\s*<div[^>]+counterDescription">problems\s*<br[^>]*>\s*solved for all time`,
		`(?is)(\d+)\s*problems.{1,150}solved for all time`,
	}
	for _, p := range patterns {
		re := regexp.MustCompile(p)
		match := re.FindStringSubmatch(html)
		if len(match) > 1 {
			val, _ := strconv.Atoi(match[1])
			return val
		}
	}
	return 0
}

func (s *SyncEngine) fetchStatusStats(u *models.User) (count int, pts float64, hard int, week int, rating int, apiHidden int, peak int, err error) {
	resp, err := http.Get("https://codeforces.com/api/user.status?handle=" + u.Handle)
	if err != nil {
		return count, pts, hard, week, 0, apiHidden, peak, fmt.Errorf("CF user.info failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("CF status error: %d", resp.StatusCode)
	}

	var res struct {
		Status string         `json:"status"`
		Result []CFSubmission `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return 0, 0, 0, 0, 0, 0, 0, err
	}

	if res.Status != "OK" {
		return 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("CF API status: %s", res.Status)
	}

	seasonStart := time.Unix(s.startDate, 0)
	daysSinceStart := int(time.Since(seasonStart).Hours() / 24)
	if daysSinceStart < 0 {
		daysSinceStart = 0
	}
	currentWeek := daysSinceStart / 7
	currentWeekStartUnix := seasonStart.AddDate(0, 0, currentWeek*7).Unix()
	problemMap := make(map[string]struct {
		solved bool
		fails  int
	})
	prevSolved := make(map[string]bool)

	for _, sub := range res.Result {
		id := fmt.Sprintf("%d%s", sub.Problem.ContestId, sub.Problem.Index)
		if sub.Verdict == "OK" && sub.CreationTimeSeconds < s.startDate {
			prevSolved[id] = true
		}
	}

	for i := len(res.Result) - 1; i >= 0; i-- {
		sub := res.Result[i]
		if sub.CreationTimeSeconds < s.startDate {
			continue
		}
		id := fmt.Sprintf("%d%s", sub.Problem.ContestId, sub.Problem.Index)
		if prevSolved[id] {
			continue
		}

		p := problemMap[id]
		if sub.Verdict == "OK" {
			if !p.solved {
				p.solved = true
				count++
				if sub.Problem.Rating >= 800 {
					pts += float64(sub.Problem.Rating-700) / 100.0
				} else if sub.Problem.Rating == 0 {
					pts += 4.0
					apiHidden++
				} else {
					pts += 1.0
				}
				if sub.CreationTimeSeconds >= currentWeekStartUnix {
					week++
					if sub.Problem.Rating > peak {
						peak = sub.Problem.Rating
					}
				}
				if p.fails >= 3 {
					hard++
				}
			}
		} else if sub.Verdict != "TESTING" && sub.Verdict != "COMPILATION_ERROR" {
			if !p.solved {
				p.fails++
			}
		}
		problemMap[id] = p
	}

	infoResp, err := http.Get("https://codeforces.com/api/user.info?handles=" + u.Handle)
	if err == nil {
		defer infoResp.Body.Close()
		var infoRes struct {
			Status string                 `json:"status"`
			Result []struct{ Rating int } `json:"result"`
		}
		json.NewDecoder(infoResp.Body).Decode(&infoRes)
		if len(infoRes.Result) > 0 {
			rating = infoRes.Result[0].Rating
		}
	}

	return count, pts, hard, week, rating, apiHidden, peak, nil
}

func (s *SyncEngine) getRankTier(pts float64) string {
	switch {
	case pts < 50:
		return "كحيان"
	case pts < 150:
		return "روش"
	case pts < 250:
		return "باشا ستراكشر"
	case pts < 350:
		return "شكسبير"
	case pts < 450:
		return "تنين مجنح"
	case pts < 600:
		return "The GOAT"
	default:
		return "CP MASTER"
	}
}

func (s *SyncEngine) fetchAtCoderStats(handle string) (float64, error) {
	if handle == "" {
		return 0, nil
	}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).Unix()
	url := fmt.Sprintf("https://kenkoooo.com/atcoder/atcoder-api/v3/user/submissions?user=%s&from_second=%d", handle, startOfMonth)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("atcoder api status: %d", resp.StatusCode)
	}

	var submissions []struct {
		ProblemID string `json:"problem_id"`
		Result    string `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&submissions); err != nil {
		return 0, err
	}

	uniqueAC := make(map[string]bool)
	for _, sub := range submissions {
		if sub.Result == "AC" {
			uniqueAC[sub.ProblemID] = true
		}
	}
	return float64(len(uniqueAC) * 5), nil
}
