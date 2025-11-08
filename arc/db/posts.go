package db

import (
	"database/sql"
	"strconv"
	"time"
)

type Post struct {
	Id        string
	CreatedAt string
	UpdatedAt string
	UserId    string
	Post      string
}

func CreatePostsTable(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY UNIQUE,
		created_at DATE not null,
		updated_at DATE not null,
		user_id INTEGER NOT NULL REFERENCES users(id),
		post TEXT NOT NULL,
	);
`

	return db.Exec(sql)
}

func AddPost(userId string, post string) *Post {
	date := time.Now().Format("2025-01-02")
	db := ConnectToDB()

	defer db.Close()

	row, err := db.Exec(`INSERT INTO posts (created_at, updated_at, user_id, post) VALUES (?,?,?,?) RETURNING *;`, date, date, userId, post)
	id, err := row.LastInsertId()

	if err != nil {

	}

	return &Post{
		Id:        strconv.Itoa(int(id)),
		CreatedAt: date,
		UpdatedAt: date,
		UserId:    userId,
		Post:      post,
	}
}

func GetAllPosts() []Post {
	db := ConnectToDB()
	rows, _ := db.Query(`Select id, created_at, updated_at, user_id, post FROM posts;`)

	var posts []Post

	for rows.Next() {

		post := &Post{}
		rows.Scan(&post.Id, &post.CreatedAt, &post.UpdatedAt, &post.UserId, &post.Post)
		posts = append(posts, *post)
	}

	return posts

}

func GetPost() {}

func UpdatePost() {}

func DeletePost() {}

func GetPostsForUser(userId string) []Post {
	db := ConnectToDB()

	defer db.Close()

	rows, _ := db.Query(`Select id, created_at, updated_at, user_id, post FROM posts WHERE posts.user_id==?;`, userId)

	var posts []Post

	for rows.Next() {

		post := &Post{}
		rows.Scan(&post.Id, &post.CreatedAt, &post.UpdatedAt, &post.UserId, &post.Post)
		posts = append(posts, *post)
	}

	return posts
}
