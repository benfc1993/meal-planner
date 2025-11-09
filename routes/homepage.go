package routes

import (
	"errors"
	"html/template"
	"log"
	"math/rand/v2"
	"meal-choices/db"
	"meal-choices/db/schema"
	"meal-choices/db/tables"
	"net/http"
	"strconv"
	"time"
)

type HomepageData struct {
	Recipes []schema.Recipe
	Week    string
}

func HandleHomepage(w http.ResponseWriter, r *http.Request, templates *template.Template) error {

	weekDate, err := getStartOfWeek(time.Now())

	if err != nil {
		return nil
	}

	data := &HomepageData{}
	data.Week = weekDate

	recipes, err := tables.GetRecipesForWeek(weekDate)

	if err == nil {
		data.Recipes = recipes
	}

	return templates.ExecuteTemplate(w, "/", data)
}

func HandleWeek(w http.ResponseWriter, r *http.Request, templates *template.Template) error {

	date, err := time.Parse(time.DateOnly, r.PathValue("week"))

	weekDate, err := getStartOfWeek(date)

	if err != nil {
		w.WriteHeader(404)
	}

	data := &HomepageData{}
	data.Week = weekDate

	recipes, err := tables.GetRecipesForWeek(weekDate)

	if err == nil {
		data.Recipes = recipes
	}

	return templates.ExecuteTemplate(w, "/week", data)
}

func HandleRecipesGenerate(w http.ResponseWriter, r *http.Request, templates *template.Template) error {
	weekDate := r.URL.Query().Get("week")
	r.ParseForm()
	count, _ := strconv.Atoi(r.Form.Get("count"))

	data := &HomepageData{}
	recentRecipes, err := tables.GetRecentRecipes()
	var recentIds = []string{}

	for i := range len(recentRecipes) {
		recentIds = append(recentIds, strconv.Itoa(recentRecipes[i].Id))
	}

	allRecipes, err := tables.GetRecipesExcept(recentIds)

	if err != nil || count > len(allRecipes) {
		return templates.ExecuteTemplate(w, "/", data)
	}

	for i := len(allRecipes) - 1; i > 0; i-- {
		j := rand.IntN(i + 1)
		allRecipes[i], allRecipes[j] = allRecipes[j], allRecipes[i]
	}

	recipes := allRecipes[:count]

	data.Recipes = recipes
	data.Week = weekDate

	return templates.ExecuteTemplate(w, "results", data)

}

func HandleAddWeekRecipes(w http.ResponseWriter, r *http.Request, templates *template.Template) error {
	r.ParseForm()
	weekDate := r.Form.Get("week")
	var ids []int
	for key, value := range r.Form {
		if key == "recipe" {

			for i := range value {
				id, err := strconv.Atoi(value[i])
				if err != nil {
					log.Printf("Attempted to insert invalid recipe id: %v for week %v\n", value[i], weekDate)
					continue
				}
				ids = append(ids, id)
			}
		}
	}

	db := db.ConnectToDB()
	db.Exec(`INSERT INTO week_recipes ()`)

	tables.InsertRecipesForWeek(weekDate, ids)

	data := &HomepageData{}
	data.Week = weekDate

	recipes, err := tables.GetRecipesForWeek(weekDate)

	if err == nil {
		data.Recipes = recipes
	}

	return templates.ExecuteTemplate(w, "/", data)
}

var days = [7]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

func getStartOfWeek(date time.Time) (string, error) {

	day := date.Format("Mon")

	index := -1
	for i := range len(days) {
		if days[i] == day {
			index = i
			break
		}
	}

	if index == -1 {
		return "", errors.New("Invalid date")
	}

	return date.AddDate(0, 0, -index).Format(time.DateOnly), nil

}
