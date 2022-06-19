package middleware

import "net/http"

type MiddlewareFunc func(http.Handler) http.Handler

// Middleware interface is anything which implements a MiddlewareFunc named Middleware.
type Middleware interface {
	Middleware(handler http.Handler) http.Handler
}

// Middleware allows MiddlewareFunc to implement the middleware interface.
func (mw MiddlewareFunc) Middleware(handler http.Handler) http.Handler {
	return mw(handler)
}
