package badges

type Badge struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Awarded     bool   `json:"awarded"`
}