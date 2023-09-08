package database

import (
	"database/sql"
)

var db *sql.DB

func Init(connStr string) (*sql.DB, error) {

	//connStr = "postgresql://nithin:office@localhost/Todo_App?sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

//func InitializeSchema(db *sql.DB) error {
//	// Execute SQL statements to create tables, indexes, etc.
//	_, err := db.Exec(`
//		CREATE TABLE IF NOT EXISTS tasks (
//			id SERIAL PRIMARY KEY,
//			name VARCHAR(255) NOT NULL,
//			description TEXT,
//			due_date DATE,
//			status VARCHAR(20),
//			created_at TIMESTAMP,
//			updated_at TIMESTAMP,
//			deleted_at TIMESTAMP,
//			user_id INT
//		);
//
//		CREATE TABLE IF NOT EXISTS users (
//			id SERIAL PRIMARY KEY,
//			username VARCHAR(15) NOT NULL,
//			password VARCHAR(255) NOT NULL
//		);
//		-- Add any other tables and schema definitions as needed.
//	`)
//	if err != nil {
//		return fmt.Errorf("failed to initialize schema: %v", err)
//	}
//
//	return nil
//}
