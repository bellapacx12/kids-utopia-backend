package rules

type Engine struct {
	rules []Rule
}

// NewEngine creates a new rules engine
func NewEngine(rules []Rule) *Engine {
	return &Engine{
		rules: rules,
	}
}

// Process runs all matching rules and returns rewards
func (e *Engine) Process(event Event, state State) ([]Reward, error) {

	var rewards []Reward

	for _, rule := range e.rules {

		if !rule.Match(event, state) {
			continue
		}

		r, err := rule.Execute(event, state)
		if err != nil {
			return nil, err
		}

		rewards = append(rewards, r...)
	}

	return rewards, nil
}