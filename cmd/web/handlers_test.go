package main

import (
	"bytes"
	"github.com/golangcollege/sessions"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"rahulrai.in/snippetbox/pkg/models/mock"
	"testing"
	"time"
)

func TestPing(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ping(rr, r)

	rs := rr.Result()

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()

	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "Ok" {
		t.Errorf("want body to equal %q", "Ok")
	}
}

func TestPingIntegration(t *testing.T) {
	// Create a new instance of our application struct. For now, this just
	//contains a couple of mock loggers (which discard anything written to
	//them).
	app := &application{
		errorLog: log.New(ioutil.Discard, "", 0),
		infoLog:  log.New(ioutil.Discard, "", 0),
		snippets: &mock.SnippetModel{},
		users:    &mock.UserModel{},
	}

	// We then use the httptest.NewTLSServer() function to create a new test
	//server, passing in the value returned by our app.routes() method as the
	//handler for the server. This starts up a HTTPS server which listens on a
	//randomly-chosen port of your local machine for the duration of the test.
	//Notice that we defer a call to ts.Close() to shutdown the server when
	//the test finishes.
	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	// The network address that the test server is listening on is contained
	//in the ts.URL field. We can use this along with the ts.Client().Get()
	//method to make a GET /ping request against the test server. This
	//returns a http.Response struct containing the response.
	rs, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}

	// We can then check the value of the response status code and body using
	//the same code as before.
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "Ok" {
		t.Errorf("want body to equal %q", "Ok")
	}
}

func TestShowSnippet(t *testing.T) {
	templateCache, err := newTemplateCache("./../../ui/html/")
	if err != nil {
		t.Fatal(err)
	}

	session := sessions.New([]byte("3dSm5MnygFHh7XidAtbskXrjbwfoJcbJ"))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	app := &application{
		errorLog:      log.New(ioutil.Discard, "", 0),
		infoLog:       log.New(ioutil.Discard, "", 0),
		snippets:      &mock.SnippetModel{},
		users:         &mock.UserModel{},
		session:       session,
		templateCache: templateCache,
	}

	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("An old silent pond...")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippet/1.23", http.StatusNotFound, nil},
		{"String ID", "/snippet/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs, err := ts.Client().Get(ts.URL + tt.urlPath)
			if err != nil {
				t.Fatal(err)
			}

			defer rs.Body.Close()
			body, err := ioutil.ReadAll(rs.Body)
			if err != nil {
				t.Fatal(err)
			}

			if rs.StatusCode != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, rs.StatusCode)
			}
			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}
}
