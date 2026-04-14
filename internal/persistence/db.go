package persistence

import "database/sql"

func openDBConnection(tableName string) (*sql.DB, error) {
	return sql.Open("superbot", tableName)
}
