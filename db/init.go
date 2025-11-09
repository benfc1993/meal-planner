package db

import (
	"database/sql"
	"log"
	"meal-choices/db/schema"

	_ "modernc.org/sqlite"
)

func ConnectToDB() *sql.DB {
	path := "./my.db"

	db, err := sql.Open("sqlite", path)

	if err != nil {
		log.Fatal("Could not connect to DB")
	}

	return db

}

func Init() {

	db := ConnectToDB()
	defer db.Close()

	schema.CreateRecipeTable(db)
	schema.CreateWeeksTable(db)
}
