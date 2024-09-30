package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Post struct {
	Id          int
	Title       string
	Description string
	Content     string
}

func main() {
	dbURL := os.Getenv("TURSO_DATABASE_URL")
	dbToken := os.Getenv("TURSO_AUTH_TOKEN")
	dbUrl := fmt.Sprintf("%s?authToken=%s", dbURL, dbToken)

	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", dbUrl, err)
		os.Exit(1)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to query: %s", err)
		os.Exit(1)
	}

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.Id, &post.Title, &post.Description, &post.Content); err != nil {
			fmt.Fprintf(os.Stderr, "failed to scan: %s", err)
			os.Exit(1)
		}
		fmt.Println(post)
	}
	defer rows.Close()

}
