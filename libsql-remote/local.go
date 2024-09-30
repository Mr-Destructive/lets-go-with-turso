package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/go-libsql"
)

func main() {
	dbName := "file:./local.db"

	db, err := sql.Open("libsql", dbName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s", err)
		os.Exit(1)
	}
	defer db.Close()
	rows, err := db.Query("SELECT ABS(RANDOM()%5) FROM generate_series(1,10)")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to query %s", err)
		os.Exit(1)
	}
	defer rows.Close()
	for rows.Next() {
		var i int
		if err := rows.Scan(&i); err != nil {
			fmt.Fprintf(os.Stderr, "failed to scan %s", err)
			os.Exit(1)
		}
		fmt.Println(i)
	}

}
