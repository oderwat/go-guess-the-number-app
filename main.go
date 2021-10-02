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

	updateAvailable bool
	url             *url.URL
}

type guessTheNumber struct {
	app.Compo

	myNumber int
	message  string // why does this not work anymore when this is "Message" instead?
	guess    string
}

var _ app.Mounter = (*guessTheNumber)(nil) // Verify the implementation
var _ app.Navigator = (*mainPage)(nil)     // Verify the implementation

func (p *mainPage) OnNav(ctx app.Context) {
	p.url = ctx.Page().URL()
}

func (c *mainPage) Render() app.UI {
	return app.Main().Body(&guessTheNumber{},
		app.Hr(),
		app.If(c.updateAvailable,
			app.Button().
				Text("Update!").
				OnClick(c.onUpdateClick),
		))
}

func (a *mainPage) onUpdateClick(ctx app.Context, e app.Event) {
	// Reloads the page to display the modifications.
	ctx.Reload()
}

// OnAppUpdate satisfies the app.AppUpdater interface. It is called when the app
// is updated in background.
func (a *mainPage) OnAppUpdate(ctx app.Context) {
	a.updateAvailable = ctx.AppUpdateAvailable() // Reports that an app update is available.
}

func (g *guessTheNumber) OnMount(ctx app.Context) {
	g.myNumber = rand.Intn(100) + 1
	g.message = "I think of a number between 1 and 100!"
	g.guess = ""

	el := app.Window().GetElementByID("guess")
	el.Call("focus")
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
				ID("guess").
				Type("text").
				Value(g.guess).
				Placeholder("Your guess (1-100)?").
				//AutoFocus(true). // does not work anyway
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
