package tables

import (
	"database/sql"
	"log"
	"meal-choices/db"
	"meal-choices/db/schema"
)

func GetRecipeCount() int {
	db := db.ConnectToDB()

	defer db.Close()

	rows, err := db.Query(`SELECT COUNT(*) FROM recipes`)

	if err != nil {
		log.Println("Error fetching recipe count")
		return 0
	}

	var count int
	for rows.Next() {
		rows.Scan(&count)
	}

	return count

}

func GetAllRecipes() ([]schema.Recipe, error) {
	db := db.ConnectToDB()

	defer db.Close()

	rows, err := db.Query(`SELECT * FROM recipes;`)

	if err != nil {
		log.Println("Error fetching all recipes")
		return nil, err
	}

	var recipes []schema.Recipe
	for rows.Next() {
		r := &schema.Recipe{}
		rows.Scan(&r.Id, &r.Name, &r.Book, &r.Page)
		recipes = append(recipes, *r)
	}

	return recipes, nil

}

func GetRecipesExcept(exclude []string) ([]schema.Recipe, error) {
	db := db.ConnectToDB()

	defer db.Close()

	query := `SELECT * FROM recipes WHERE id NOT IN(`

	for i := range len(exclude) {
		query += exclude[i] + ","
	}

	if len(exclude) > 0 {
		query = query[:len(query)-1]
	}

	query += ");"

	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	var recipes []schema.Recipe
	for rows.Next() {
		r := &schema.Recipe{}
		rows.Scan(&r.Id, &r.Name, &r.Book, &r.Page)
		recipes = append(recipes, *r)
	}

	return recipes, nil

}

func GetRecipeById(id int) (*schema.Recipe, error) {
	db := db.ConnectToDB()

	defer db.Close()

	rows, err := db.Query(`SELECT id, name, book, page FROM recipes WHERE id = ?`, id)

	if err != nil {
		log.Println("sad")
		return nil, err
	}

	defer rows.Close()

	recipe := &schema.Recipe{}

	rows.Next()
	err = rows.Scan(&recipe.Id, &recipe.Name, &recipe.Book, &recipe.Page)

	if err != nil {
		log.Println("sad")
		return nil, err
	}
	return recipe, nil

}

func AddRecipe(name string, book string, page int) (sql.Result, error) {
	db := db.ConnectToDB()

	defer db.Close()

	return db.Exec(`INSERT INTO recipes (name, book, page) values (?,?,?);`, name, book, page)

}
