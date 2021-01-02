package view

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

var (
	LayoutPath = "views/layouts"
	LayoutName = "application"
)

func Render(name string, data interface{}) ([]byte, error) {
	parts := strings.Split(name, "/")
	var namespace string
	var action string

	if len(parts) == 2 {
		namespace, action = parts[0], parts[1]
	} else {
		namespace, action = "", parts[0]
	}

	var content bytes.Buffer
	contentTemplate := template.Must(template.New(fmt.Sprintf("%s.gohtml", action)).ParseGlob(fmt.Sprintf("views/%s/*.gohtml", namespace)))
	err := contentTemplate.Execute(&content, data)
	if err != nil {
		return nil, err
	}

	funcMap := map[string]interface{}{
		"yield": func() string { return content.String() },
	}

	layout := template.Must(template.New(fmt.Sprintf("%s.gohtml", LayoutName)).Funcs(funcMap).ParseFiles(fmt.Sprintf("%s/%s.gohtml", LayoutPath, LayoutName)))
	var buf bytes.Buffer
	err = layout.Execute(&buf, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
