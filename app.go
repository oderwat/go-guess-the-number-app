package main

import (
	"net/url"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type appPage struct {
	app.Compo
	autoUpdate      bool
	updateAvailable bool
	url             *url.URL
}

var _ app.Mounter = (*guessTheNumber)(nil) // Verify the implementation
var _ app.Navigator = (*appPage)(nil)      // Verify the implementation

func (a *appPage) OnNav(ctx app.Context) {
	if app.IsServer {
		return
	}
	a.url = ctx.Page().URL()
	app.Logf(a.url.String())
	if a.url.Fragment == "autoupdate" {
		a.autoUpdate = true
	} else {
		a.autoUpdate = false
	}
}

func (a *appPage) Render() app.UI {
	count := 1 // try 1000 if you like
	body := make([]app.UI, 0, count+5)
	body = append(body,
		app.If(a.updateAvailable,
			func() app.UI {
				return app.Button().
					Text("Update!").
					OnClick(a.onUpdateClick)
			}),
		app.Hr())
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

func (a *appPage) onUpdateClick(ctx app.Context, _ app.Event) {
	// Reloads the page to display the modifications.
	ctx.Reload()
}

// OnAppUpdate satisfies the app.AppUpdater interface. It is called when the app
// is updated in background.
func (a *appPage) OnAppUpdate(ctx app.Context) {
	a.updateAvailable = ctx.AppUpdateAvailable() // Reports that an app update is available.
	if a.updateAvailable && a.autoUpdate {
		ctx.Reload()
	}
}
