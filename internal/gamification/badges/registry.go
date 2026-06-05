package badges

var Registry = map[string]Badge{
	"first_page": {
		ID:          "first_page",
		Title:       "🎉 First page read!",
		Description: "You read your first page!",
		Icon:        "book-open",
	},

	"book_completed": {
		ID:          "book_completed",
		Title:       "📚 Book completed!",
		Description: "You finished reading a book!",
		Icon:        "trophy",
	},

	"streak_7": {
		ID:          "streak_7",
		Title:       "🔥 7-day streak!",
		Description: "You read for 7 days in a row!",
		Icon:        "fire",
	},
}
func GetBadge(id string) (Badge, bool) {
	b, ok := Registry[id]
	return b, ok
}