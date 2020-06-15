package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable, app.noSurf, app.authenticate)

	// build routes using pat. Pat will create a mux rather than http.NewServeMux below.
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	//mux := http.NewServeMux()
	//mux.HandleFunc("/", app.home)
	////mux.HandleFunc("/snippet", app.showSnippet)
	////mux.HandleFunc("/snippet/create", app.createSnippet)

	// Auth routes
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	// Ping
	mux.Get("/ping", http.HandlerFunc(ping))

	fileServer := http.FileServer(http.Dir("../../ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	//mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// secure header middleware wrapping mux middleware wrapping handlers in handlers.go
	//return app.recoverPanic(app.logRequest(secureHeaders(mux)))
	return standardMiddleware.Then(mux)
}
