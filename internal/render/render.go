package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
)

func check(fMessage string, err error) {
	if err != nil {
		log.Fatalf(fMessage, err)
	}
}

func IncludeTemplate(path string, data map[string]any) http.HandlerFunc {
	const basePath = "./web/templates/"

	return http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		tLayout, err := template.ParseFiles(basePath + "layout/main.html")
		check("Error while parsing layout file: %v", err)

		buffer := new(bytes.Buffer)
		tPage, err := template.ParseFiles(basePath + path)
		check("Error while parsing page file: %v", err)

		err = tPage.Execute(buffer, data)
		check("Error while executing page template: %v", err)

		layoutData := struct {
			Content template.HTML
			Title   any
		}{
			Content: template.HTML(buffer.String()),
			Title:   data["title"],
		}

		err = tLayout.Execute(res, layoutData)
		check("Error while executing layout template: %v", err)
	})
}
