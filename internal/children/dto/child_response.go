package dto

import "time"

type ChildResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Age       *int    `json:"age,omitempty"`
	Language  string  `json:"language"`
    Gamification *GamificationDTO `json:"gamification,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
type GamificationDTO struct {
	XP     int         `json:"xp"`
	Level  int         `json:"level"`
	Streak int         `json:"streak"`
	Badges []BadgeDTO  `json:"badges"`
	Milestones []MilestoneDTO `json:"milestones"`
	Themes []ThemeDTO  `json:"themes"`
}
type BadgeDTO struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Awarded     bool   `json:"awarded"`
}
type MilestoneDTO struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`

	Current int  `json:"current"`
	Target  int  `json:"target"`
    ProgressPercent int    `json:"progress_percent"`
	Awarded bool `json:"awarded"`
}
type ThemeDTO struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	Unlocked  bool   `json:"unlocked"`
}