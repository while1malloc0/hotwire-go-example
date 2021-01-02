package controllers

import (
	"html/template"

	unrolledRender "github.com/unrolled/render"

	"github.com/while1malloc0/hotwire-go-example/pkg/timefmt"
)

var render = unrolledRender.New(
	unrolledRender.Options{
		Extensions: []string{".gohtml"},
		Layout:     "layout",
		Funcs:      []template.FuncMap{timefmt.FuncMap},
	},
)
