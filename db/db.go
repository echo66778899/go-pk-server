package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLManager struct {
	db *sql.DB
}

func NewMySQLManager(username, password, host, port, database string) (*MySQLManager, error) {
	// Create the MySQL connection string
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)

	// Open a connection to the MySQL database
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Return the MySQLManager instance
	return &MySQLManager{
		db: db,
	}, nil
}

func (m *MySQLManager) Close() error {
	// Close the database connection
	err := m.db.Close()
	if err != nil {
		return err
	}

	return nil
}

func (m *MySQLManager) Execute(query string, args ...interface{}) (sql.Result, error) {
	// Execute the SQL query
	result, err := m.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *MySQLManager) Query(query string, args ...interface{}) (*sql.Rows, error) {
	// Execute the SQL query and return the result set
	rows, err := m.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
