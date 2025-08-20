package middleware

import "net/http"

type Middleware func(http.Handler) (http.Handler, error)

type Chain struct {
	constructors []Middleware
}

func New(constructors ...Middleware) Chain {
	return Chain{constructors: constructors}
}

func (c Chain) Then(h http.Handler) (http.Handler, error) {
	var err error
	for i := range c.constructors {
		h, err = c.constructors[len(c.constructors)-1-i](h)
		if err != nil {
			return nil, err
		}
	}
	return h, nil
}

func (c Chain) Append(constructors ...Middleware) Chain {
	newCons := make([]Middleware, len(c.constructors)+len(constructors))
	copy(newCons, c.constructors)
	copy(newCons[len(c.constructors):], constructors)
	return Chain{constructors: newCons}
}
