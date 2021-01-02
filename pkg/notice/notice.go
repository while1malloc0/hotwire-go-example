package notice

import (
	"context"
	"net/http"
)

type contextKey struct{}

var ContextKey = contextKey{}

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

func Set(w http.ResponseWriter, value string) {
	http.SetCookie(w, &http.Cookie{Name: "notice", Value: value, Path: "/"})
}

func Clear(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: "notice", MaxAge: -1})
}

func Get(r *http.Request) string {
	var notice string
	if noticeCookie, err := r.Cookie("notice"); err == nil {
		notice = noticeCookie.Value
	}
	return notice
}
