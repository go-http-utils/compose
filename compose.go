package compose

import (
	"fmt"
	"net/http"
	"reflect"
)

type middleware struct {
	fun  interface{}
	opts []interface{}
}

type Composer struct {
	h           http.Handler
	middlewares []middleware
}

func New(h http.Handler) *Composer {
	return &Composer{h: h, middlewares: []middleware{}}
}

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

func (c Composer) Handler() http.Handler {
	return c.h
}
