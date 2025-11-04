package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

type Recipe struct {
	Id   int
	Name string
	Book string
	Page int
}

func NewRecipe() *Recipe {
	return &Recipe{
		Id:   -1,
		Name: "",
		Book: "",
		Page: -1,
	}
}

func InitDb() (*sql.DB, error) {
	db, err := ConnectToDB()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer db.Close()

	_, err = createTables(db)

	return db, err
}

func ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./my.db")

	return db, err

}

func createTables(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS recipes (
		id INTEGER PRIMARY KEY UNIQUE,
		name TEXT NOT NULL,
		book TEXT NOT NULL,
		page INTEGER DEFAULT 0,
		UNIQUE(name, book)
	);
	CREATE TABLE IF NOT EXISTS week_recipes (
		date DATE not null,
		recipe_id INTEGER NOT NULL REFERENCES recipes(id)
	);
	`

	return db.Exec(sql)
}

func GetRecipeCount() int {
	db, err := ConnectToDB()

	if err != nil {
		return 0
	}

	defer db.Close()

	rows, err := db.Query(`SELECT COUNT(*) FROM recipes`)

	var count int
	for rows.Next() {
		rows.Scan(&count)
	}

	return count

}

func GetAllRecipes() ([]Recipe, error) {
	db, err := ConnectToDB()

	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(`SELECT * FROM recipes;`)

	var recipes []Recipe
	for rows.Next() {
		r := &Recipe{}
		rows.Scan(&r.Id, &r.Name, &r.Book, &r.Page)
		recipes = append(recipes, *r)
	}

	return recipes, nil

}

func GetRecipeById(id int) (*Recipe, error) {
	db, err := ConnectToDB()

	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(`SELECT id, name, book, page FROM recipes WHERE id = ?`, id)

	if err != nil {
		log.Println("sad")
		return nil, err
	}

	defer rows.Close()

	recipe := &Recipe{}

	rows.Next()
	err = rows.Scan(&recipe.Id, &recipe.Name, &recipe.Book, &recipe.Page)

	if err != nil {
		log.Println("sad")
		return nil, err
	}
	return recipe, nil

}

func AddRecipe(name string, book string, page int) (sql.Result, error) {
	db, err := ConnectToDB()

	if err != nil {
		return nil, err
	}

	defer db.Close()

	return db.Exec(`INSERT INTO recipes (name, book, page) values (?,?,?);`, name, book, page)

}

func GetRecipesForWeek(date string) ([]Recipe, error) {

	db, err := ConnectToDB()

	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(`SELECT recipes.id, recipes.name, recipes.book, recipes.page FROM recipes 
	INNER JOIN week_recipes 
	ON recipes.id=week_recipes.recipe_id AND week_recipes.date=?;`, date)

	if err != nil {
		log.Println("Oops")
		return nil, err
	}

	defer rows.Close()

	var recipes []Recipe

	for rows.Next() {
		println("recipe")
		r := &Recipe{}
		rows.Scan(&r.Id, &r.Name, &r.Book, &r.Page)
		recipes = append(recipes, *r)
	}

	return recipes, nil
}

func InsertRecipesForWeek(date string, recipes []int) error {
	db, err := ConnectToDB()

	if err != nil {
		return err
	}

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
