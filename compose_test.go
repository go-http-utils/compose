package compose

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ComposeSuite struct {
	suite.Suite

	m1    bool
	m1t   time.Time
	m2    bool
	m2t   time.Time
	h     bool
	req   *http.Request
	res   http.ResponseWriter
	start time.Time
}

func (s *ComposeSuite) SetupTest() {
	s.start = time.Now()
	s.req = httptest.NewRequest(http.MethodGet, "/", nil)
	s.res = testResponseWriter{}

	s.m1, s.m2 = false, false
	s.m1t, s.m2t = s.start, s.start
}

func (s *ComposeSuite) TestOneMiddleware() {
	New(http.HandlerFunc(s.handler)).Use(
		middleware1, s,
	).Handler().ServeHTTP(s.res, s.req)

	s.True(s.h)
	s.True(s.m1)
	s.False(s.m2)
	s.True(s.m1t.After(s.start))
	s.True(s.m2t.Equal(s.start))
}

func (s *ComposeSuite) TestManyMiddlewares() {
	New(http.HandlerFunc(s.handler)).Use(
		middleware1, s,
	).Use(
		middleware2, s,
	).Handler().ServeHTTP(s.res, s.req)

	s.True(s.h)
	s.True(s.m1)
	s.True(s.m2)
	s.True(s.m1t.After(s.start))
	s.True(s.m2t.After(s.start))
}

func (s *ComposeSuite) TestUseMiddlewares() {
	New(http.HandlerFunc(s.handler)).UseMiddlewares(
		[]Middleware{
			Middleware{Fun: middleware1, Opts: []interface{}{s}},
			Middleware{Fun: middleware2, Opts: []interface{}{s}},
		},
	).Handler().ServeHTTP(s.res, s.req)

	s.True(s.h)
	s.True(s.m1)
	s.True(s.m2)
	s.True(s.m1t.After(s.start))
	s.True(s.m2t.After(s.start))
}

func (s *ComposeSuite) TestArgumentsPanic() {
	s.Panics(func() {
		New(http.HandlerFunc(s.handler)).Use(
			func(s string) {}, s,
		).Handler().ServeHTTP(s.res, s.req)
	})
}

func (s *ComposeSuite) TestReturnValuePanic() {
	s.Panics(func() {
		New(http.HandlerFunc(s.handler)).Use(
			func(h http.Handler) {}, s,
		).Handler().ServeHTTP(s.res, s.req)
	})
}

func (s *ComposeSuite) handler(res http.ResponseWriter, req *http.Request) {
	s.h = true
}

func TestCompose(t *testing.T) {
	suite.Run(t, new(ComposeSuite))
}

func middleware1(h http.Handler, cs *ComposeSuite) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cs.m1 = true
		cs.m1t = time.Now()
		h.ServeHTTP(res, req)
	})
}

func middleware2(h http.Handler, cs *ComposeSuite) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cs.m2 = true
		cs.m2t = time.Now()
		h.ServeHTTP(res, req)
	})
}

type testResponseWriter struct{}

func (t testResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (t testResponseWriter) WriteHeader(s int) {}

func (t testResponseWriter) Write(b []byte) (int, error) {
	return 0, nil
}
