package routes

import (
	"fmt"
	"html/template"
	"log"
	"meal-choices/db/schema"
	"meal-choices/db/tables"
	"net/http"
	"strconv"
	"strings"
)

func HandleRecipeAdd(w http.ResponseWriter, r *http.Request, templates *template.Template) error {

	r.ParseForm()
	name := r.Form.Get("name")
	book := r.Form.Get("book")
	pageNum, _ := strconv.Atoi(r.Form.Get("page"))

	if name == "" || book == "" {
		var missingValue []string
		if book == "" {
			missingValue = append(missingValue, "book")
		}
		if name == "" {
			missingValue = append(missingValue, "name")
		}

		w.WriteHeader(422)
		templates.ExecuteTemplate(w, "recipe-error", fmt.Sprintf("Problem creating recipe, missing: %v", strings.Join(missingValue, ", ")))
		return templates.ExecuteTemplate(w, "recipe-form", &schema.Recipe{Id: -1, Name: name, Book: book, Page: pageNum})
	}

	_, err := tables.AddRecipe(name, book, pageNum)
	if err != nil {
		message := "Problem creating recipe"
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			message = fmt.Sprintf("Recipe: \"%v\" in Book: \"%v\" already exists.", name, book)
		}

		w.WriteHeader(422)
		templates.ExecuteTemplate(w, "recipe-error", message)
		return templates.ExecuteTemplate(w, "recipe-form", &schema.Recipe{Id: -1, Name: name, Book: book, Page: pageNum})
	}

	w.WriteHeader(201)
	templates.ExecuteTemplate(w, "recipe-result", fmt.Sprintf("Recipe \"%v\" added.", name))
	return templates.ExecuteTemplate(w, "recipe-form", &schema.Recipe{})
}

func HandleGetAllRecipes(w http.ResponseWriter, r *http.Request, templates *template.Template) error {
	recipes, err := tables.GetAllRecipes()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(500)
		return templates.ExecuteTemplate(w, "recipes-list", nil)
	}

	return templates.ExecuteTemplate(w, "recipes-list", recipes)
}
