package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid creds")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}

// This acts like constructor.
func NewSnippet(id int, title string, content string, created time.Time,
	expires time.Time) *Snippet {
	return &Snippet{
		ID:      id,
		Title:   title,
		Content: content,
		Created: created,
		Expires: expires,
	}
}
