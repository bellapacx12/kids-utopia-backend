package themes

import "context"

type Service struct {
	repo Repository
}

func New(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Unlock(ctx context.Context, childID, themeID string) error {
	return s.repo.Unlock(ctx, childID, themeID)
}

func (s *Service) GetUnlocked(ctx context.Context, childID string) ([]string, error) {
	return s.repo.GetUnlocked(ctx, childID)
}
func (s *Service) BuildUserThemes(ctx context.Context, childID string) ([]UserTheme, error) {
	unlocked, err := s.repo.GetUnlockedThemes(ctx, childID)
	if err != nil {
		return nil, err
	}

	unlockedMap := make(map[string]bool)
	for _, id := range unlocked {
		unlockedMap[id] = true
	}

	result := make([]UserTheme, 0, len(Registry))

	for _, t := range Registry {
		result = append(result, UserTheme{
			ID:     t.ID,
			Locked: !unlockedMap[t.ID],
		})
	}

	return result, nil
}