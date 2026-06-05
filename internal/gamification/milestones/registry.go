package milestones

type Milestone struct {
	ID          string
	Title       string
	Description string
	Target      int
}

var Registry = map[string]Milestone{
	"first_page": {
		ID:          "first_page",
		Title:       "🎉 First page read!",
		Description: "Read your first page",
		Target:      1,
	},

	"book_completed": {
		ID:          "book_completed",
		Title:       "📚 Book completed!",
		Description: "Finish a book",
		Target:      1,
	},

	"streak_7": {
		ID:          "streak_7",
		Title:       "🔥 7-day streak!",
		Description: "Read for 7 days in a row",
		Target:      7,
	},
}
func Get(id string) (Milestone, bool) {
	m, ok := Registry[id]
	return m, ok
}