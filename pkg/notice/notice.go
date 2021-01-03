// Package notice implements cookie-based flash/notice functionality for
// persisting ephemeral state between requests
package notice

import (
	"context"
	"net/http"
)

type contextKey struct{}

// ContextKey is the type-safe key for storing a notice
var ContextKey = contextKey{}

// Context is an http.Handler that parses a notice from a request's cookies. If
// one is found, it's added to the request's Context and removed from cookies
func Context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notice := Get(r)
		if notice != "" {
			// Notice has been collected, so expire it now
			Clear(w)
			ctx := context.WithValue(r.Context(), ContextKey, notice)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

// Set updates the notice on a given response
func Set(w http.ResponseWriter, value string) {
	http.SetCookie(w, &http.Cookie{Name: "notice", Value: value, Path: "/"})
}

// Clear removes a notice on a given response
func Clear(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: "notice", MaxAge: -1})
}

// Get retrieves a notice on a given response
func Get(r *http.Request) string {
	var notice string
	if noticeCookie, err := r.Cookie("notice"); err == nil {
		notice = noticeCookie.Value
	}
	return notice
}
