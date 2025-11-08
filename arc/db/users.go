package db

import "database/sql"

func CreateUserTable(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY UNIQUE,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		UNIQUE(name, email)
	);
	`
	return db.Exec(sql)
}

func AddUser()    {}
func DeleteUser() {}
