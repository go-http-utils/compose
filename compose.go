package compose

import (
	"fmt"
	"net/http"
	"reflect"
)

// Version is this package's version
const Version = "0.1.0"

// Middleware includes the middleware function and the
// arguments to be passed.
type Middleware struct {
	// Func is the middleware function. The first argument of the
	// function should be a http.Handler and the return value should also be a
	// http.Handler ( middlewareFunc(h http.Handler, ...) http.Handler ).
	Fun interface{}
	// Opts are the arguments to pass to the middleware function.
	Opts []interface{}
}

// Composer wraps the inner http.Handler and can be used to
// append middlewares to it.
type Composer struct {
	h http.Handler
}

// New returns a new Composer instance.
func New(h http.Handler) *Composer {
	return &Composer{h}
}

// Use appends the middleware to the inner handler. Accept the middleware
// function as the first argument and the first argument of the middleware
// function should be a http.Handler and the return value should also be
// http.Handler ( middlewareFunc(h http.Handler, ...) http.Handler ).
func (c *Composer) Use(fun interface{}, opts ...interface{}) *Composer {
	t := reflect.TypeOf(fun)

	if t.NumIn() == 0 || t.In(0).Name() != "Handler" {
		panic(fmt.Sprintf("compose: middleware function: %s"+
			" should take http.Handler as the first argument.", t))
	}

	if t.NumOut() != 1 || t.Out(0).Name() != "Handler" {
		panic(fmt.Sprintf("compose: middleware function: %s"+
			" should return http.Handler as the only return value.", t))
	}

	optVals := make([]reflect.Value, len(opts)+1)
	optVals[0] = reflect.ValueOf(c.h)

	for i, opt := range opts {
		optVals[i+1] = reflect.ValueOf(opt)
	}

	c.h = reflect.ValueOf(fun).Call(optVals)[0].Interface().(http.Handler)

	return c
}

// UseMiddlewares appends the specfied middlewares to the inner handler.
func (c *Composer) UseMiddlewares(ms []Middleware) *Composer {
	for _, m := range ms {
		c.Use(m.Fun, m.Opts...)
	}

	return c
}

// Handler returns the inner handler.
func (c Composer) Handler() http.Handler {
	return c.h
}
