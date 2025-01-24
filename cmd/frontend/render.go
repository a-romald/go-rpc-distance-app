package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
)

//go:embed templates/*
var templateFS embed.FS

var loader = NewLoader("templates", templateFS)
var views = jet.NewSet(loader)

// Render func renders the page using Jet templates
func Render(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println("Unexpected template err:", err.Error())
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println("Error executing template:", err)
		return err
	}
	return nil
}
