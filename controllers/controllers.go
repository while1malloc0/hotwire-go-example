package controllers

import (
	"html/template"

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
