package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tursodatabase/go-libsql"
	_ "github.com/tursodatabase/go-libsql"
)

func main() {

	dbName := "local.db"
	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, dbName)
	fmt.Println(dbPath)

	dbURL := os.Getenv("TURSO_DATABASE_URL")
	dbAuthToken := os.Getenv("TURSO_AUTH_TOKEN")

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, dbURL, libsql.WithAuthToken(dbAuthToken))
	fmt.Println(connector)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s", err)
		os.Exit(1)
	}

	defer connector.Close()

	db := sql.OpenDB(connector)
	fmt.Println("Connected to database")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s", err)
		os.Exit(1)
	}
	defer db.Close()

	createPostTableQuery := `CREATE TABLE IF NOT EXISTS posts(
        id INTEGER PRIMARY KEY,
        title VARCHAR(100),
        description VARCHAR(255),
        content TEXT
    );`

	_, err = db.Exec(createPostTableQuery)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create table %s", err)
		os.Exit(1)
	}
	createPostQuery := `INSERT INTO posts(title, description, content) 
        VALUES(?, ?, ?)`

	_, err = db.Exec(createPostQuery, "test title", "test description", "test content")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to insert %s", err)
		os.Exit(1)
	}

	_, err = connector.Sync()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to sync %s", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully synced %s db\n", dbPath)
	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to query %s", err)
		os.Exit(1)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var title string
		var description string
		var content string
		if err := rows.Scan(&id, &title, &description, &content); err != nil {
			fmt.Fprintf(os.Stderr, "failed to scan %s", err)
			os.Exit(1)
		}
		fmt.Println(id, title, description, content)
	}

}
