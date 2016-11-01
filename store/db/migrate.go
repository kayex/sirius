package db

import "fmt"

func (db *Db) CreateSchema() {
	queries := []string{
		createUsersTable(),
	}

	db.execMany(queries)
}

func (db *Db) execMany(queries []string) {
	for _, q := range queries {
		_, err := db.conn.Exec(q)

		if err != nil {
			fmt.Printf("Executing migration query failed: %v\n", q)
			panic(err)
		}
	}
}

func createUsersTable() string {
	return `CREATE TABLE users (id VARCHAR(255) NOT NULL, name VARCHAR(255), PRIMARY KEY(id)`;
}
