package db

import (
	"database/sql"
	"fmt"
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

	_, err = createTable(db)

	return db, err
}

func ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./my.db")

	return db, err

}

func createTable(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS recipes (
		id INTEGER PRIMARY KEY UNIQUE,
		name TEXT NOT NULL,
		book TEXT NOT NULL,
		page INTEGER DEFAULT 0,
		UNIQUE(name, book)
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

	rows, err := db.Query(`SELECT * FROM recipes`)

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
		return nil, err
	}

	defer rows.Close()

	recipe := &Recipe{}

	err = rows.Scan(&recipe.Id, &recipe.Name, &recipe.Book, &recipe.Page)

	if err != nil {
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
