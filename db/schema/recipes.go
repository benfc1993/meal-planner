package schema

import "database/sql"

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

func CreateRecipeTable(db *sql.DB) (sql.Result, error) {
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
