package main

import (
	"math/rand"
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type guessTheNumber struct {
	app.Compo

	myNumber int
	message  string // why does this not work anymore when this is "Message" instead?
	guess    string
	guesses  []string
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
