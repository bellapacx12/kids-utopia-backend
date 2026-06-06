package service

import (
	"context"
	"log"

	dto "github.com/bellapacx/kids-utopia/internal/children/dto"
	"github.com/bellapacx/kids-utopia/internal/gamification/badges"
	milestones "github.com/bellapacx/kids-utopia/internal/gamification/milestones"
	milestoneservice "github.com/bellapacx/kids-utopia/internal/gamification/milestones"
	"github.com/bellapacx/kids-utopia/internal/gamification/model"
	"github.com/bellapacx/kids-utopia/internal/gamification/registry"
	"github.com/bellapacx/kids-utopia/internal/gamification/repository"
	"github.com/bellapacx/kids-utopia/internal/gamification/rules"
	themes "github.com/bellapacx/kids-utopia/internal/gamification/themes"
	progressservice "github.com/bellapacx/kids-utopia/internal/progress/service"
	streakservice "github.com/bellapacx/kids-utopia/internal/streak/service"
)

type Service struct {
	repo   repository.Repository
	engine *rules.Engine
	milestoneService *milestoneservice.Service
	streakService   *streakservice.StreakService
	progress *progressservice.ProgressService
	themeService     *themes.Service
}
type Snapshot struct {
	XP     int
	Level  int
	Streak int
	Badges []badges.Badge
	Milestones []dto.MilestoneDTO
	Themes     []dto.ThemeDTO
}
type Milestone struct {
	ID          string
	Title       string
	Description string
	Current     int
	Target      int
	ProgressPercent int    `json:"progress_percent"`
	Awarded     bool
}
type ThemeDTO struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Icon     string `json:"icon"`
	Unlocked bool   `json:"unlocked"`
}

func New(repo repository.Repository,ms *milestoneservice.Service, streak *streakservice.StreakService, progress *progressservice.ProgressService, themeService *themes.Service) *Service {
	return &Service{
		repo:   repo,
		engine: rules.NewEngine(registry.NewRules()),
		milestoneService: ms,
		streakService: streak,
		progress: progress,
		themeService: themeService,
	}
}
func (s *Service) ProcessEvent(ctx context.Context, event rules.Event) error {

	

	state, err := s.repo.GetState(ctx, event.ChildID)
if err != nil {
	log.Printf("❌ STATE LOAD FAILED child=%s err=%v", event.ChildID, err)
	return err
}

seen, err := s.progress.Exists(ctx, event.ChildID, event.BookID)
if err != nil {
	log.Printf("❌ PROGRESS CHECK FAILED child=%s book=%s err=%v",
		event.ChildID, event.BookID, err)
	return err
}

state.BookSeen = seen
	

	rewards, err := s.engine.Process(event, state)
	if err != nil {
		log.Printf("❌ RULE ENGINE FAILED: %v", err)
		return err
	}

	log.Printf("🎁 REWARDS GENERATED: count=%d", len(rewards))

	for i, r := range rewards {
		log.Printf("➡️ reward[%d]: type=%s value=%d meta=%v",
			i, r.Type, r.Value, r.Meta,
		)
	}


	for _, reward := range rewards {

		switch reward.Type {

		// =========================
		// XP SYSTEM
		// =========================
		case rules.RewardTypeXP:

			log.Printf("⭐ XP reward: %+v", reward)

			err := s.repo.InsertTransaction(
				ctx,
				&model.XPTransaction{
					ChildID:  event.ChildID,
					Source:   event.Type,
					SourceID: event.EventID,
					XPAmount: reward.Value,
				},
			)
			if err != nil {
				log.Printf("❌ XP TRANSACTION FAILED: %v", err)
				return err
			}

			err = s.repo.UpsertXP(ctx, event.ChildID, reward.Value)
			if err != nil {
				log.Printf("❌ XP UPDATE FAILED: %v", err)
				return err
			}


		// =========================
		// BADGES / STREAKS
		// =========================
		case rules.RewardTypeBadge, rules.RewardTypeStreak:

			log.Printf("🏅 MILESTONE reward: %+v", reward)

			err := s.milestoneService.Process(ctx, event, rewards)
			if err != nil {
				log.Printf("❌ MILESTONE PROCESS FAILED: %v", err)
				return err
			}

			log.Println("✅ MILESTONE APPLIED")

		// =========================
		// THEMES
		// =========================
		case rules.RewardTypeTheme:

			log.Printf("🎨 THEME reward: %+v", reward)

			raw, ok := reward.Meta["theme_id"]
			if !ok {
				log.Printf("❌ THEME META MISSING: %+v", reward.Meta)
				continue
			}

			themeID, ok := raw.(string)
			if !ok {
				log.Printf("❌ THEME META INVALID TYPE: %T", raw)
				continue
			}

			log.Printf("🔓 unlocking theme: %s for child=%s", themeID, event.ChildID)

			err := s.themeService.Unlock(ctx, event.ChildID, themeID)
			if err != nil {
				log.Printf("❌ THEME UNLOCK FAILED: %v", err)
				return err
			}

			log.Println("✅ THEME UNLOCKED")

		default:
			log.Printf("⚠️ UNKNOWN REWARD TYPE: %s", reward.Type)
		}
	}

	log.Println("🎉 GAMIFICATION COMPLETE")
	return nil
}
func (s *Service) GetXP(
	ctx context.Context,
	childID string,
) (*model.ChildXP, error) {
	return s.repo.GetXP(ctx, childID)
}
func (s *Service) GetSnapshot(ctx context.Context, childID string) (*Snapshot, error) {

	// =========================
	// 1. XP
	// =========================
	xpData, err := s.repo.GetXP(ctx, childID)
	if err != nil {
		log.Println("XP ERROR:", err)
		return nil, err
	}

	xp := 0
	if xpData != nil {
		xp = xpData.TotalXP
	}

	level := calculateLevel(xp)

	// =========================
	// 2. STREAK
	// =========================
	streak := 0
	streakData, err := s.streakService.GetStreak(ctx, childID)
	if err != nil {
		log.Printf("streak not found for child %s: %v", childID, err)
	} else if streakData != nil {
		streak = streakData.CurrentStreak
	}

	// =========================
	// 3. BADGES (awarded + full registry)
	// =========================
	rawBadges, err := s.repo.GetBadges(ctx, childID)
	if err != nil {
		log.Println("badges ERROR:", err)
		return nil, err
	}

	awarded := make(map[string]bool)
	for _, b := range rawBadges {
		awarded[b.MilestoneID] = true
	}

	badgesResult := make([]badges.Badge, 0, len(badges.Registry))
	for _, badge := range badges.Registry {
		badgesResult = append(badgesResult, badges.Badge{
			ID:          badge.ID,
			Title:       badge.Title,
			Description: badge.Description,
			Icon:        badge.Icon,
			Awarded:     awarded[badge.ID],
		})
	}

	// =========================
	// 4. PROGRESS (needed for milestones)
	// =========================
	progressList, err := s.progress.ListByChild(ctx, childID)
	if err != nil {
		return nil, err
	}

	// =========================
	// 5. AWARDED MILESTONES (same table as badges)
	// =========================
	awardedMilestones := make(map[string]bool)
	for _, b := range rawBadges {
		awardedMilestones[b.MilestoneID] = true
	}

	// =========================
	// 6. BUILD MILESTONES
	// =========================
	milestonesResult := make([]dto.MilestoneDTO, 0)

	// ---- FIRST PAGE ----
	if m, ok := milestones.Get("first_page"); ok {

		current := 0
		for _, p := range progressList {
			if p.CurrentPage > 0 {
				current = 1
				break
			}
		}

		milestonesResult = append(milestonesResult, dto.MilestoneDTO{
			ID:          m.ID,
			Title:       m.Title,
			Description: m.Description,
			Current:     current,
			Target:      m.Target,
			ProgressPercent: progressPercent(current, m.Target),
			Awarded:     awardedMilestones["first_page"],
		})
	}

	// ---- BOOK COMPLETED ----
	if m, ok := milestones.Get("book_completed"); ok {

		completed := 0
		for _, p := range progressList {
			if p.Completed {
				completed++
			}
		}

		milestonesResult = append(milestonesResult, dto.MilestoneDTO{
			ID:          m.ID,
			Title:       m.Title,
			Description: m.Description,
			Current:     completed,
			Target:      m.Target,
			ProgressPercent: progressPercent(completed, m.Target),
			Awarded:     awardedMilestones["book_completed"],
		})
	}

	// ---- STREAK 7 ----
	if m, ok := milestones.Get("streak_7"); ok {

		milestonesResult = append(milestonesResult, dto.MilestoneDTO{
			ID:          m.ID,
			Title:       m.Title,
			Description: m.Description,
			Current:     streak,
			Target:      m.Target,
			ProgressPercent: progressPercent(streak, m.Target),
			Awarded:     awardedMilestones["streak_7"],
		})
	}
    unlockedList, err := s.themeService.GetUnlocked(ctx, childID)
if err != nil {
	unlockedList = []string{}
}

unlocked := map[string]bool{}
for _, id := range unlockedList {
	unlocked[id] = true
}

theme := make([]dto.ThemeDTO, 0, len(themes.Registry))

for _, t := range themes.Registry {
	theme = append(theme, dto.ThemeDTO{
		ID:       t.ID,
		Name:    t.Name,
		Icon:     t.Icon,
		Unlocked: unlocked[t.ID],
	})
}
	// =========================
	// FINAL RESPONSE
	// =========================
	return &Snapshot{
		XP:         xp,
		Level:      level,
		Streak:     streak,
		Badges:     badgesResult,
		Milestones: milestonesResult,
		Themes: theme,
	}, nil
}
func calculateLevel(xp int) int {
	switch {
	case xp < 100:
		return 1
	case xp < 300:
		return 2
	case xp < 600:
		return 3
	default:
		return 4
	}
}
func progressPercent(current, target int) int {
	if target <= 0 {
		return 0
	}

	p := (current * 100) / target

	if p > 100 {
		return 100
	}

	return p
}