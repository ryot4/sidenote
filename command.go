package main

type Command interface {
	Name() string
	Description() string
	Run(args []string, options *Options) error
}

type SyntaxError struct {
	message string
}

func NewSyntaxError(message string) error {
	return &SyntaxError{message: message}
}

func (e *SyntaxError) Error() string {
	return e.message
}

var (
	ErrTooManyArgs   = NewSyntaxError("too many arguments")
	ErrNoFileName    = NewSyntaxError("no file specified")
	ErrNoDstFileName = NewSyntaxError("no destination file specified")
)
