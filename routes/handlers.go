package routes

import (
	"embed"
	"html/template"
	"net/http"
)

func Handlers(f embed.FS, templates *template.Template) {
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, f, r.URL.Path)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { HandleHomepage(w, templates) })
}
