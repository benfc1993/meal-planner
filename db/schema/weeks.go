package schema

import "database/sql"

func CreateWeeksTable(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS week_recipes (
		date DATE not null,
		recipe_id INTEGER NOT NULL REFERENCES recipes(id)
	);
	`

	return db.Exec(sql)
}
