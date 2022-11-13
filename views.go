package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

//generatePage:
func generatePage(writer http.ResponseWriter, name string, data interface{}) {
	files := []string{"layout", "user.sidebar", strings.Join([]string{"content", name}, ".")}
	t := parseTemplateFiles("layout", files...)
	t.ExecuteTemplate(writer, "layout", data)
}

//generateLoginPage
func generateLoginPage(writer http.ResponseWriter, message interface{}) {
	t := parseTemplateFiles("login", "layout.login")
	t.Execute(writer, message)
}

//parseTemplateFiles this function is used to parse all templates and
// return it as one template type, and you have to remember that location of templates is templates/index.html
func parseTemplateFiles(layout string, filenames ...string) (t *template.Template) {
	var files []string
	t = template.New(layout)
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	t = template.Must(t.ParseFiles(files...))
	return
}

//generateJsonResponse
func generateJSON(writer http.ResponseWriter, body interface{}) error {
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(body)
	return err
}
