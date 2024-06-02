package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func main() {
	rand.Seed(time.Now().UnixNano()) // needs to be better for a real app

	// Components routing:
	app.Route("/", app.NewZeroComponentFactory(&appPage{}))

	app.RunWhenOnBrowser()

	if !app.IsServer {
		return // let go dead code elimination optimize the frontend code (saves 11 mb in app.wasm)
	}

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
