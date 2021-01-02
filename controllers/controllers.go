package controllers

import (
	"context"
	"html/template"
	"net/http"

	"github.com/while1malloc0/hotwire-go-example/pkg/timefmt"

	unrolledRender "github.com/unrolled/render"
)

var render = unrolledRender.New(
	unrolledRender.Options{
		Extensions: []string{".gohtml"},
		Layout:     "layout",
		Funcs:      []template.FuncMap{timefmt.FuncMap},
	},
)

var ContextKeyNotice = contextKey{}

func NoticeContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notice := GetNotice(r)
		if notice != "" {
			// Notice has been collected, so expire it now
			ClearNotice(w)
			ctx := context.WithValue(r.Context(), ContextKeyNotice, notice)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

func SetNotice(w http.ResponseWriter, value string) {
	http.SetCookie(w, &http.Cookie{Name: "notice", Value: value, Path: "/"})
}

func ClearNotice(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: "notice", MaxAge: -1})
}

func GetNotice(r *http.Request) string {
	var notice string
	if noticeCookie, err := r.Cookie("notice"); err == nil {
		notice = noticeCookie.Value
	}
	return notice
}
