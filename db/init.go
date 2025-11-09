package db

import (
	"database/sql"
	"log"
	"meal-choices/db/schema"
	"os"

	_ "modernc.org/sqlite"
)

func ConnectToDB() *sql.DB {
	path := "./my.db"

	if os.Getenv("ENV") != "dev" {
		dbDirPath := os.Getenv("HOME") + "/.config/testing"
		_, err := os.ReadDir(dbDirPath)
		if err != nil {
			os.Mkdir(os.Getenv("HOME")+"/.config/testing", 0777)
		}
		path = os.Getenv("HOME") + "/.config/testing/my.db"
	}

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
