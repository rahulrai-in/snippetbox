package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"rahulrai.in/snippetbox/pkg/models"
	"time"

	// Only init will be invoked. A dot removes the need of qualifier.
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"rahulrai.in/snippetbox/pkg/models/mysql"
)

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

// DI
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	// change this to inline interface to improve testability
	//snippets      *mysql.SnippetModel
	snippets interface {
		Insert(string, string, string) (int, error)
		Get(int) (*models.Snippet, error)
		Latest() ([]*models.Snippet, error)
	}
	templateCache map[string]*template.Template
	session       *sessions.Session
	//users         *mysql.UserModel
	users interface {
		Insert(string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	}
}

func init() {
	fmt.Println("Welcome to snippetbox")
}

func main() {
	// var mux *http.ServeMux = http.NewServeMux()
	addr := flag.String("addr", ":4000", "Server port")
	dsn := flag.String("dsn", "root:password@tcp(localhost:3306)/snippetbox?parseTime=true", "MySQL data source name")

	// session secret
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret Key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INF\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("../../ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteStrictMode

	app := application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		session:       session,
		users:         &mysql.UserModel{DB: db},
	}

	mux := app.routes()

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      mux,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute, // == 1 * time.Minute
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Println("Starting server on port " + *addr)
	err = srv.ListenAndServeTLS("../../tls/cert.pem", "../../tls/key.pem")
	errorLog.Fatal(err)

	// Handle errors only in main
	// Loggers created from log.New() are type safe so create instances of logger within goroutines
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
