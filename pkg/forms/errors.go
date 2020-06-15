package forms

type errors []string

func (e *errors) Add(msg string) {
	*e = append(*e, msg)
}

func (e *errors) Get() []string {
	return *e
}
