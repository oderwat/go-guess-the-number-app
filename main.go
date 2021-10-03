package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func main() {
	rand.Seed(time.Now().UnixNano()) // needs to be better for a real app

	// Components routing:
	app.Route("/", &appPage{})

	app.RunWhenOnBrowser()

	// HTTP routing:
	http.Handle("/", &app.Handler{
		Name:        "Guess the Number",
		Description: "A simple game using go-app as framework!",
		Scripts: []string{
			"https://livejs.com/live.js", // Add simple live reloading
		},
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
