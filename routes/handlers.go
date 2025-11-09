package routes

import (
	"embed"
	"html/template"
	"net/http"
)

func Handlers(f embed.FS, templates *template.Template) {
	http.HandleFunc("GET /static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, f, r.URL.Path)
	})
	http.HandleFunc("/", handler(HandleHomepage, templates))
	http.HandleFunc("POST /recipes/generate", handler(HandleRecipesGenerate, templates))
	http.HandleFunc("POST /recipes/add", handler(HandleHomepage, templates))
	http.HandleFunc("GET /all", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "/all", nil)
	})
	http.HandleFunc("GET /recipes/all", handler(HandleGetAllRecipes, templates))
	http.HandleFunc("DELETE /", func(w http.ResponseWriter, r *http.Request) { w.Write(nil) })
}

func handler(fn func(w http.ResponseWriter, r *http.Request, templates *template.Template) error, templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, templates)
	}
}
