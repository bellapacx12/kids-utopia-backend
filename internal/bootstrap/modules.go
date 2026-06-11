package bootstrap

import (
	"github.com/bellapacx/kids-utopia/internal/modules"
)

func (a *App) registerModules() {

	modules.RegisterAuth(
		a.Router,
		a.Container,
	)

	modules.RegisterUsers(
		a.Router,
		a.Container,
	)

	modules.RegisterChildren(
		a.Router,
		a.Container,
	)

	modules.RegisterProgress(
		a.Router,
		a.Container,
	)

	modules.RegisterSubscriptions(
		a.Router,
		a.Container,
	)

	modules.RegisterBooks(
		a.Router,
		a.Container,
	)
	modules.RegisterReaderSession(
	a.Router,
	a.Container,
) 
modules.RegisterReader(
	a.Router,
	a.Container,
) 
   modules.RegisterBookmarks(
	a.Router,
	a.Container,
   )
   modules.RegisterRecommendation(
	a.Router,
	a.Container,
   )
}