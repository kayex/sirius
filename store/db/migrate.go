package db

func (db *Db) CreateSchema() {
}

func createUsersTable() string {
	return `CREATE TABLE users (id VARCHAR(255) NOT NULL, name VARCHAR(255), PRIMARY KEY(id)`;
}
