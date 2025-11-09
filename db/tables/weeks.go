package tables

import (
	"context"
	"log"
	"meal-choices/db"
	"meal-choices/db/schema"
	"time"
)

func GetRecentRecipes() ([]schema.Recipe, error) {
	db := db.ConnectToDB()

	defer db.Close()

	date := time.Now().AddDate(0, -1, 0).Format("2006-01-02")

	rows, err := db.Query(`SELECT recipes.id, recipes.name, recipes.book, recipes.page 
	FROM recipes 
	INNER JOIN week_recipes 
	ON recipes.id=week_recipes.recipe_id AND date>?;`, date)

	if err != nil {
		log.Println("Oops")
		return nil, err
	}

	defer rows.Close()

	var recipes []schema.Recipe

	for rows.Next() {
		println("recipe")
		r := &schema.Recipe{}
		rows.Scan(&r.Id, &r.Name, &r.Book, &r.Page)
		recipes = append(recipes, *r)
	}

	return recipes, nil
}

func GetRecipesForWeek(date string) ([]schema.Recipe, error) {

	db := db.ConnectToDB()

	defer db.Close()

	rows, err := db.Query(`SELECT recipes.id, recipes.name, recipes.book, recipes.page FROM recipes 
	INNER JOIN week_recipes 
	ON recipes.id=week_recipes.recipe_id AND week_recipes.date=?;`, date)

	if err != nil {
		log.Println("Oops")
		return nil, err
	}

	defer rows.Close()

	var recipes []schema.Recipe

	for rows.Next() {
		println("recipe")
		r := &schema.Recipe{}
		rows.Scan(&r.Id, &r.Name, &r.Book, &r.Page)
		recipes = append(recipes, *r)
	}

	return recipes, nil
}

func InsertRecipesForWeek(date string, recipes []int) error {
	db := db.ConnectToDB()

	defer db.Close()

	trx, err := db.BeginTx(context.Background(), nil)

	for i := range recipes {
		trx.Exec(`INSERT INTO week_recipes (date, recipe_id) values (?,?)`, date, recipes[i])
	}

	err = trx.Commit()
	if err != nil {
		return err
	}

	return nil
}
