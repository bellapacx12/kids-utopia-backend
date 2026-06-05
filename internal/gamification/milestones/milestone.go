package milestones

type MilestoneDTO struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`

	Current int  `json:"current"`
	Target  int  `json:"target"`

	Awarded bool `json:"awarded"`
}
type RewardType string

const (
	RewardBadge RewardType = "badge"
	RewardTheme RewardType = "theme"
)