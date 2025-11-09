package routes

import (
	"errors"
	"html/template"
	"log"
	"math/rand/v2"
	"meal-choices/db/schema"
	"meal-choices/db/tables"
	"net/http"
	"strconv"
	"time"
)

type HomepageData struct {
	Recipes []schema.Recipe
}

func HandleHomepage(w http.ResponseWriter, r *http.Request, templates *template.Template) error {

	weekDate, err := getStartOfWeekDate()

	if err != nil {
		return nil
	}

	data := &HomepageData{}

	recipes, err := tables.GetRecipesForWeek(weekDate)

	if err == nil {
		data.Recipes = recipes
	}

	return templates.ExecuteTemplate(w, "/", data)
}

func HandleRecipesGenerate(w http.ResponseWriter, r *http.Request, templates *template.Template) error {
	log.Println("Generate")
	r.ParseForm()
	count, _ := strconv.Atoi(r.Form.Get("count"))
	log.Println(count)

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

	return templates.ExecuteTemplate(w, "results", data)

}

var days = [7]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

func getStartOfWeekDate() (string, error) {
	today := time.Now().Format("Mon")

	index := -1
	for i := range len(days) {
		if days[i] == today {
			index = i
			break
		}
	}

	if index == -1 {
		return "", errors.New("Invalid date")
	}

	return time.Now().AddDate(0, 0, -index).Format("2006-01-02"), nil

}
