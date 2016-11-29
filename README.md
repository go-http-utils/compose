# compose
[![Build Status](https://travis-ci.org/go-http-utils/compose.svg?branch=master)](https://travis-ci.org/go-http-utils/compose)
[![Coverage Status](https://coveralls.io/repos/github/go-http-utils/compose/badge.svg?branch=master)](https://coveralls.io/github/go-http-utils/compose?branch=master)

Compose your http middlewares with happiness.

## Installation

```
go get -u github.com/go-http-utils/compose
```

## Documentation

Documentation can be found here: https://godoc.org/github.com/go-http-utils/compose

## Usage

It transforms:

```go
mux := http.NewServeMux()

mux.HandleFunc("/", middleware1(middleware2(middleware3(handler1, arg3)), arg1, arg2))
mux.HandleFunc("/example", middleware3(middleware4(handler2, arg2, arg3), arg1))

http.ListenAndServe(":8080", middleware4(middleware5(mux)))
```

to:

```go
mux := http.NewServeMux()

mux.HandleFunc("/", compose.New(handler1).
  Use(middleware1, arg1, arg2).
  Use(middleware2).
  Use(middleware3, arg3).
  Handler())

mux.HandleFunc("/example", compose.New(handler2).
  UseMiddlewares([]Middleware{
    Middleware{Func: middleware3, Opts: []interface{arg1}},
    Middleware{Func: middleware4, Opts: []interface{arg2, arg3}},
  }).
  Handler())

http.ListenAndServe(":8080", compose.New(mux).
  Use(middleware4).
  Use(middleware5).
  Handler()
)
```
