package controllers

import (
	unrolledRender "github.com/unrolled/render"
)

var render = unrolledRender.New(
	unrolledRender.Options{
		Extensions: []string{".gohtml"},
		Layout:     "layout",
	},
)
