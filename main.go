package main

import (
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type mainPage struct {
	app.Compo

	url *url.URL
}

type guessTheNumber struct {
	app.Compo

	myNumber int
	message  string
	guess    string
}

var _ app.Mounter = (*guessTheNumber)(nil) // Verify the implementation
var _ app.Navigator = (*mainPage)(nil)     // Verify the implementation

func (p *mainPage) OnNav(ctx app.Context) {
	p.url = ctx.Page().URL()
}

func (c *mainPage) Render() app.UI {
	return app.Body().Body(&guessTheNumber{})
}

func (g *guessTheNumber) OnMount(ctx app.Context) {
	g.myNumber = rand.Intn(100) + 1
	g.message = "I think of a number between 1 and 100!"
	g.guess = ""
}

func (g *guessTheNumber) onEnter(call func(ctx app.Context, e app.Event)) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		app.Logf("KeyboardEvent: %v", e.Value.Get("code"))
		if e.Value.Get("code").String() != "Enter" {
			return
		}
		call(ctx, e)
	}
}

func (g *guessTheNumber) guessEvent(ctx app.Context, e app.Event) {
	v, err := strconv.Atoi(ctx.JSSrc().Get("value").String())
	if err != nil {
		v = 0
	}
	if v < g.myNumber {
		g.message = "Your number is smaller!"
	} else if v > g.myNumber {
		g.message = "Your number is bigger!"
	} else {
		g.message = "Congratulations! You found my number!"
	}
	ctx.JSSrc().Call("select")
}

func (g *guessTheNumber) Render() app.UI {
	return app.Div().Body(
		app.H1().Body(
			app.Text("Guess my Number!"),
		),
		app.P().Body(
			app.P().Body(app.Text(g.message)),
			app.Input().
				Type("text").
				Value(g.guess).
				Placeholder("Your guess (1-100)?").
				AutoFocus(true).
				OnKeyup(g.onEnter(g.guessEvent)),
		),
	)
}

func main() {
	rand.Seed(time.Now().UnixNano()) // needs to be better for a real app

	// Components routing:
	app.Route("/", &mainPage{})
	//app.Route("/hello", &hello{})

	app.RunWhenOnBrowser()

	// HTTP routing:
	http.Handle("/", &app.Handler{
		Name:        "Guess the Number",
		Description: "A simple game!",
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
