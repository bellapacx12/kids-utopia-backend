package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bellapacx/kids-utopia/internal/children/dto"
	"github.com/bellapacx/kids-utopia/internal/children/service"
	"github.com/bellapacx/kids-utopia/internal/gamification/badges"
	themes "github.com/bellapacx/kids-utopia/internal/gamification/themes"
	"github.com/bellapacx/kids-utopia/pkg/contextkeys"
)

type ChildHandler struct {
	service *service.ChildService
}

func NewChildHandler(s *service.ChildService) *ChildHandler {
	return &ChildHandler{service: s}
}

func (h *ChildHandler) Create(c *gin.Context) {
	parentID := c.GetString(contextkeys.UserID)

	var req dto.CreateChildRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	err := h.service.Create(c.Request.Context(), parentID, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "child created"})
}

func (h *ChildHandler) MyChildren(c *gin.Context) {
	parentID := c.GetString(contextkeys.UserID)

	children, err := h.service.GetByParent(c.Request.Context(), parentID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": children})
}

func (h *ChildHandler) GetChild(c *gin.Context) {
	childID := c.Param("childId")

	child, snap, err := h.service.FindByID(c.Request.Context(), childID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var gamification *dto.GamificationDTO

	if snap != nil {
		gamification = &dto.GamificationDTO{
			XP:     snap.XP,
			Level:  snap.Level,
			Streak: snap.Streak,
			Badges: mapBadges(snap.Badges),
			Milestones: snap.Milestones,
			Themes: snap.Themes,
		}
	}

	c.JSON(http.StatusOK, dto.ChildResponse{
		ID:           child.ID,
		Name:         child.Name,
		AvatarURL:    child.AvatarURL,
		Age:          child.Age,
		Language:     child.Language,
		Gamification: gamification,
		CreatedAt:    child.CreatedAt,
	})
}
func mapBadges(in []badges.Badge) []dto.BadgeDTO {
	out := make([]dto.BadgeDTO, 0, len(in))

	for _, b := range in {
		out = append(out, dto.BadgeDTO{
			ID:          b.ID,
			Title:       b.Title,
			Description: b.Description,
			Icon:        b.Icon,
			Awarded:     true,
		})
	}

	return out
}
func mapMilestones(in []dto.MilestoneDTO) []dto.MilestoneDTO {
	out := make([]dto.MilestoneDTO, 0, len(in))

	for _, m := range in {
		out = append(out, dto.MilestoneDTO{
			ID:          m.ID,
			Title:       m.Title,
			Description: m.Description,
			Current:     m.Current,
			Target:      m.Target,
			Awarded:     m.Awarded,
		})
	}

	return out
}
func mapThemes(in []themes.Theme, unlocked map[string]bool) []dto.ThemeDTO {
	out := make([]dto.ThemeDTO, 0, len(in))

	for _, t := range in {
		out = append(out, dto.ThemeDTO{
			ID:     t.ID,
			Name:   t.Name,
			Icon:   t.Icon,
			Unlocked: !unlocked[t.ID],
		})
	}

	return out
}