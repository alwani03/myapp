package main

import "database/sql"

// _ensureSQLImported prevents build errors by ensuring database/sql is imported.
// It is not used at runtime.
func _ensureSQLImported(_ *sql.DB) {}