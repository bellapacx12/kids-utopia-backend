package main

import (
	"log"

	"github.com/bellapacx/kids-utopia/internal/bootstrap"
)

func main() {

	app := bootstrap.NewApp()

	port := app.Container.Config.AppPort

	log.Printf("🚀 Server running on %s", port)

	if err := app.Router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
