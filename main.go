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
	autoUpdate      bool
	updateAvailable bool
	url             *url.URL
}

type guessTheNumber struct {
	app.Compo

	myNumber int
	message  string // why does this not work anymore when this is "Message" instead?
	guess    string
	guesses  []string
}

var _ app.Mounter = (*guessTheNumber)(nil) // Verify the implementation
var _ app.Navigator = (*mainPage)(nil)     // Verify the implementation

func (a *mainPage) OnNav(ctx app.Context) {
	a.url = ctx.Page().URL()
	//app.Logf(p.url.Fragment)
	if a.url.Fragment == "autoupdate" {
		app.Logf("auto update enabled")
		a.autoUpdate = true
	} else {
		app.Logf("auto update disabled")
		a.autoUpdate = false
	}
}

func (a *mainPage) Render() app.UI {
	count := 1 // try 1000 if you like
	body := make([]app.UI, 0, count+5)
	body = append(body,
		app.If(a.updateAvailable,
			app.Button().
				Text("Update!").
				OnClick(a.onUpdateClick),
			app.Hr()))
	for i := 0; i < count; i++ {
		body = append(body, &guessTheNumber{})
	}
	var aumodeFragment, aumodeText string
	if a.autoUpdate {
		aumodeFragment = "#"
		aumodeText = "Automatic updates"
	} else {
		aumodeFragment = "#autoupdate"
		aumodeText = "Manual updates"
	}
	return app.Main().Body(append(body, app.Hr(),
		app.A().Href(aumodeFragment).Text(aumodeText))...)
}

func (a *mainPage) onUpdateClick(ctx app.Context, _ app.Event) {
	// Reloads the page to display the modifications.
	ctx.Reload()
}

// OnAppUpdate satisfies the app.AppUpdater interface. It is called when the app
// is updated in background.
func (a *mainPage) OnAppUpdate(ctx app.Context) {
	a.updateAvailable = ctx.AppUpdateAvailable() // Reports that an app update is available.
	if a.updateAvailable && a.autoUpdate {
		ctx.Reload()
	}
}

func (g *guessTheNumber) OnMount(_ app.Context) {
	g.gameInit()
}

func (g *guessTheNumber) gameInit() {
	g.myNumber = rand.Intn(100) + 1
	g.message = "I think of a number between 1 and 100!"
	g.guesses = g.guesses[:0] // keep memory
	g.guess = ""

	el := app.Window().GetElementByID("guess")
	el.Call("focus")
}

func (g *guessTheNumber) onEnter(call func(ctx app.Context, e app.Event)) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		//app.Logf("KeyboardEvent: %v", e.Value.Get("code"))
		if e.Value.Get("code").String() != "Enter" {
			return
		}
		call(ctx, e)
	}
}

func (g *guessTheNumber) guessEvent(ctx app.Context, _ app.Event) {
	guess := ctx.JSSrc().Get("value").String()
	if guess == "Restart!" {
		g.gameInit()
		return
	}
	v, err := strconv.Atoi(guess)
	if err != nil {
		v = 0
	}
	if v < g.myNumber {
		g.message = "Your number is smaller!"
		g.guesses = append(g.guesses, guess+" was smaller")
	} else if v > g.myNumber {
		g.message = "Your number is bigger!"
		g.guesses = append(g.guesses, guess+" was bigger")
	} else {
		g.guesses = append(g.guesses, guess+" is correct!")
		g.message = "Congratulations! You found my number!"
		g.guess = "Restart!"
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
			app.Ul().Body(app.Range(g.guesses).Slice(func(i int) app.UI {
				return app.Li().Text(g.guesses[i])
			})),
		),
	)
}

func main() {
	rand.Seed(time.Now().UnixNano()) // needs to be better for a real app

	// Components routing:
	app.Route("/", &mainPage{})

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
