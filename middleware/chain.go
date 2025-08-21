package middleware

import (
	"github.com/nprimo/quick/web"
)

type Middleware func(next web.HandlerFuncWithError) web.HandlerFuncWithError

type Chain struct {
	constructors []Middleware
}

func New(constructors ...Middleware) Chain {
	return Chain{constructors: constructors}
}

func (c Chain) Then(h web.HandlerFuncWithError) web.HandlerFuncWithError {
	for i := range c.constructors {
		h = c.constructors[len(c.constructors)-1-i](h)
	}
	return h
}

func (c Chain) Append(constructors ...Middleware) Chain {
	newCons := make([]Middleware, len(c.constructors)+len(constructors))
	copy(newCons, c.constructors)
	copy(newCons[len(c.constructors):], constructors)
	return Chain{constructors: newCons}
}
