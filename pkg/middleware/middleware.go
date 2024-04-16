package middleware

import (
	"net/http"
)

type Middleware func(next http.Handler) http.Handler

func MiddlewareStack(stack ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(stack); i >= 0; i-- {
			mw := stack[i]
			next = mw(next)
		}
		return next
	}
}
