package main

import (
	"embed"
	"fmt"
	"html/template"
	"strings"

	"github.com/gin-gonic/gin"
)

type templateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CRSFToken string
	Flash     string
	Warning   string
	Error     string
	IsAuthed  int
	API       string
}

var functions = template.FuncMap{}

var templateFS embed.FS

func (app *application) addDefaultData(td *templateData, c *gin.Context) *templateData {
	// Modify this method to add default data to the template data based on the Gin context.
	return td
}

func (app *application) renderTemplate(c *gin.Context, page string, td *templateData, partials ...string) error {
	var t *template.Template
	var err error

	templateToRender := fmt.Sprintf("templates/%s.page.tmpl", page)
	_, templateInMap := app.templateCache[templateToRender]
	if app.config.env == "production" && templateInMap {
		t = app.templateCache[templateToRender]
	} else {
		t, err = app.parseTemplate(partials, page, templateToRender)
		if err != nil {
			app.errorlog.Println(err)
			return err
		}
	}

	if td == nil {
		td = &templateData{}
	}
	td = app.addDefaultData(td, c)

	// Use c.HTML method to render the template using Gin's context.
	err = t.Execute(c.Writer, td)
	if err != nil {
		app.errorlog.Println(err)
		return err
	}

	return nil
}

func (app *application) parseTemplate(partials []string, page, templateToRender string) (*template.Template, error) {
	var t *template.Template
	var err error

	if len(partials) > 0 {
		for i, x := range partials {
			partials[i] = fmt.Sprintf("templates/%s.partial.tmpl", x)
		}
	}
	if len(partials) > 0 {
		t, err = template.New(fmt.Sprintf("%s.page.tmpl", page)).Funcs(functions).ParseFS(
			templateFS,
			"templates/base.layout.tmpl",
			strings.Join(partials, ","),
			templateToRender,
		)
	} else {
		t, err = template.New(fmt.Sprintf("%s.page.tmpl", page)).Funcs(functions).ParseFS(
			templateFS,
			"templates/base.layout.tmpl",
			templateToRender,
		)
	}
	if err != nil {
		app.errorlog.Println(err)
		return nil, err
	}

	app.templateCache[templateToRender] = t
	return t, nil
}
