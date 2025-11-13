package routes

import (
	"embed"
	"html/template"
	"meal-choices/db/schema"
	"net/http"
)

func Handlers(f embed.FS, templates *template.Template) {

	var handler = func(fn func(w http.ResponseWriter, r *http.Request, templates *template.Template) error) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fn(w, r, templates)
		}
	}
	http.HandleFunc("GET /static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, f, r.URL.Path)
	})

	http.HandleFunc("/", handler(HandleHomepage))
	http.HandleFunc("DELETE /", func(w http.ResponseWriter, r *http.Request) { w.Write(nil) })
	http.HandleFunc("POST /recipes/generate", handler(HandleRecipesGenerate))

	http.HandleFunc("GET /week/{week}", handler(HandleWeek))
	http.HandleFunc("POST /week", handler(HandleAddWeekRecipes))

	http.HandleFunc("GET /all", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "/all", nil)
	})
	http.HandleFunc("GET /recipes/all", handler(HandleGetAllRecipes))

	http.HandleFunc("POST /recipes/add", handler(HandleRecipeAdd))
	http.HandleFunc("POST /recipes/upload", handler(HandleCsvUpload))
	http.HandleFunc("GET /add", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "/add", schema.NewRecipe())
	})

}
