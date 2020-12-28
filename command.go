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

var subCommands = []Command{
	&CatCommand{},
	&EditCommand{},
	&ImportCommand{},
	&InitCommand{},
	&LsCommand{},
	&PathCommand{},
	&RmCommand{},
	&ServeCommand{},
	&ShowCommand{},
}
