package closer

import "gopasskeeper/internal/logger"

var globalCloser = New()

// Add is a function for calling global closer Add method.
func Add(f ...func() error) {
	globalCloser.Add(f...)
}

// CloseAll is a function for calling global closer CloseAll method.
func CloseAll() {
	globalCloser.CloseAll()
}

// Closer is a structure to aggregate closer functions.
type Closer struct {
	funcs []func() error
}

// New is a builder function for Closer structure.
func New() *Closer {
	return &Closer{
		funcs: make([]func() error, 0),
	}
}

// Add is a method to add closer functions to Closer.
func (c *Closer) Add(f ...func() error) {
	c.funcs = append(c.funcs, f...)
}

// CloseAll is a method to trigger all closer functions one by one.
func (c *Closer) CloseAll() {
	for _, f := range c.funcs {
		if err := f(); err != nil {
			logger.Error("error close", err)
		}
	}
}
