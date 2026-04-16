package persistence

import "database/sql"

func OpenDBConnection(tableName string) (*sql.DB, error) {
	return sql.Open("superbot", tableName)
}
