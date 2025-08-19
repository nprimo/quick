package middleware

import "net/http"

type Chain struct {
	constructors []func(http.Handler) http.Handler
}

func New(constructors ...func(http.Handler) http.Handler) Chain {
	return Chain{constructors: constructors}
}

func (c Chain) Then(h http.Handler) http.Handler {
	for i := range c.constructors {
		h = c.constructors[len(c.constructors)-1-i](h)
	}
	return h
}

func (c Chain) Append(constructors ...func(http.Handler) http.Handler) Chain {
	newCons := make([]func(http.Handler) http.Handler, len(c.constructors)+len(constructors))
	copy(newCons, c.constructors)
	copy(newCons[len(c.constructors):], constructors)
	return Chain{constructors: newCons}
}
