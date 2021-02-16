package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

//serverError
func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

//client side error
func (app *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

//page not found
func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

//Adding current year to every page. This takes a pointer to a templateData struct
//adds, the current year to the CurrentYear field, and then return the pointer.
func (app *Application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	return td
}

//render web pages
func (app *Application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	//Retrieve the appropriate template set from the cache based on the page name
	//(like 'home.page.html'). If no entry exists in the cache with the provided name
	//call the serverError helper method
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}
	//Initialize a new buffer
	buf := new(bytes.Buffer)
	//execute template
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}
	//Write the contents of the buffer to the http.ResponseWriter.
	//Again this is another time where we pass our http.ResponseWriter
	//to a function that takes an io.Writer
	buf.WriteTo(w)
}
