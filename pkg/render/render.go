package render

import (
	"bytes"
	"fmt"
	"github.com/jjang65/go-hello-word/pkg/config"
	"github.com/jjang65/go-hello-word/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// Init functions which type is FuncMap defining the mapping from names to functions.
var functions = template.FuncMap{}

// app is the pointer to AppConfig
var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	// Init map containing string key and pointer to Template
	myCache := map[string]*template.Template{}

	// Find all pages
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// Loop through all pages and if there is any template matched,
	// return parsed layouts
	for _, page := range pages {
		name := filepath.Base(page)

		// Init templateSet containing all templates
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// Find matched layouts
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		// If any matched layout, parse all layouts
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	// Define templateCache
	var tc map[string]*template.Template

	// If UseCache is true, assign app.TemplateCache to `tc`
	if app.UseCache {
		// Get the template cache from the app config
		tc = app.TemplateCache
	} else {
		// If UseCache is false, call CreateTemplateCache() always,
		// so to create template
		tc, _ = CreateTemplateCache()
	}

	// Get template by indexing the template path
	t, ok := tc[tmpl]
	// if index tmpl does not exist, ok should be false
	if !ok {
		// Stop server
		log.Fatal("could not get template from template cache")
	}

	// Put parsed template into bytes in memory
	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	// Execute applies a parsed template to the specified data object, writing the output to wr.
	_ = t.Execute(buf, td)

	// buf.WriteTo writes data to w until the buffer is drained or an error occurs
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}
}
