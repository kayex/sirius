package db

import "fmt"

func (db *Db) CreateSchema() {
	queries := createUsersTable()
	queries = append(queries, createConfigurationsTable()...)

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

func createUsersTable() []string {
	return []string{
		`CREATE TABLE users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
			token VARCHAR(255) NOT NULL,
			created_at TIMESTAMP (0)
		)`,
		`CREATE UNIQUE INDEX uniq_token ON users USING btree (token)`,
	}
}

func createConfigurationsTable() []string {
	return []string{
		`CREATE TABLE configurations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
			user_id UUID NOT NULL,
			plugin_guid VARCHAR(255) NOT NULL,
			config JSON
		)`,
		`CREATE UNIQUE INDEX uniq_user_plugin ON configurations USING btree (user_id, plugin_guid)`,
		`ALTER TABLE configurations ADD CONSTRAINT FK_CONFIGURATION_USER FOREIGN KEY (user_id) REFERENCES users (id) NOT DEFERRABLE INITIALLY IMMEDIATE;`,
	}
}
