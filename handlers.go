package main

import (
	"net/http"
	"time"
)

//home handler
func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w) //using the notFound helper
		return
	}
	//years
	t := time.Now().Year()

	//render the home page
	app.render(w, r, "home.page.html", &templateData{
		CurrentYear: t,
	})
}
