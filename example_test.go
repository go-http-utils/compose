package compose_test

import (
	"net/http"

	"github.com/go-http-utils/compose"
)

var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

const (
	middleware1 = iota
	middleware2
	middleware3
	middleware4
	middleware5
	arg1
	arg2
	arg3
	arg4
	arg5
)

func Example() {
	mux := http.NewServeMux()

	mux.Handle("/example1", compose.New(handler).
		Use(middleware1, arg1, arg2).
		Use(middleware2, arg3).
		Use(middleware3).
		Handler())

	mux.Handle("/example2", compose.New(handler).
		UseMiddlewares([]compose.Middleware{
			compose.Middleware{Fun: middleware3, Opts: []interface{}{arg4, arg5}},
			compose.Middleware{Fun: middleware4},
		}).
		Handler())

	http.ListenAndServe(":8080", compose.New(mux).Use(middleware5).Handler())
}

func ExampleComposer_Use() {
	mux := http.NewServeMux()

	mux.Handle("/example1", compose.New(handler).
		Use(middleware1, arg1, arg2).
		Use(middleware2, arg3).
		Use(middleware3).
		Handler())

	http.ListenAndServe(":8080", compose.New(mux).Use(middleware5).Handler())
}

func ExampleComposer_UseMiddlewares() {
	mux := http.NewServeMux()

	mux.Handle("/example2", compose.New(handler).
		UseMiddlewares([]compose.Middleware{
			compose.Middleware{Fun: middleware3, Opts: []interface{}{arg4, arg5}},
			compose.Middleware{Fun: middleware4},
		}).
		Handler())

	http.ListenAndServe(":8080", compose.New(mux).Use(middleware5).Handler())
}
