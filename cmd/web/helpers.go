package main

import (
	"bytes"
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// The error did not occur here so we trace till a level up.
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("template %s does not exist", name))
		return
	}

	// write to buffer before sending to client so that any errors in HTML template don's return a partial page.
	buffer := &bytes.Buffer{}
	err := ts.Execute(buffer, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	buffer.WriteTo(w)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.CurrentYear = time.Now().Year()

	// add flash data if exists
	td.Flash = app.session.PopString(r, "flash")

	td.IsAuthenticated = app.isAuthenticated(r)

	// Add CSRF token to each page
	td.CSRFToken = nosurf.Token(r)

	return td
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(contextKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}
