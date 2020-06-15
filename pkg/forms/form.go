package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		Values: data,
		Errors: make([]string, 0),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(fmt.Sprintf("%s Value required", field))
		}
	}
}

func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(fmt.Sprintf("Field %s can not be longer than %d characters", field, d))
	}
}

func (f *Form) AllowedValues(field string, opts ...string) {
	value := f.Get(field)
	for _, opt := range opts {
		if value == opt {
			return
		}
	}

	f.Errors.Add(fmt.Sprintf("%s is invalid"))
}

func (f *Form) IsValid() bool {
	return len(f.Errors) == 0
}
