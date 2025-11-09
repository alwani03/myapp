package db

import (
    "database/sql"
    "io/fs"
    "log"
    "os"
    "path/filepath"
    "sort"
    "strings"
)

// RunMigrations executes .sql files in internal/app/db/migrations filtered by dialect.
// Naming convention:
//   *_all.sql      -> runs for any dialect
//   *_pg.sql       -> runs for Postgres
//   *_mysql.sql    -> runs for MySQL
func RunMigrations(db *sql.DB, dialect string) error {
    migrationsDir := filepath.Join("internal", "app", "db", "migrations")
    entries, err := os.ReadDir(migrationsDir)
    if err != nil {
        return err
    }
    files := make([]fs.DirEntry, 0)
    for _, e := range entries {
        if e.IsDir() || filepath.Ext(e.Name()) != ".sql" {
            continue
        }
        name := e.Name()
        if strings.HasSuffix(name, "_all.sql") ||
            (strings.EqualFold(dialect, "postgres") && strings.HasSuffix(name, "_pg.sql")) ||
            (strings.EqualFold(dialect, "mysql") && strings.HasSuffix(name, "_mysql.sql")) {
            files = append(files, e)
        }
    }
    sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
    for _, f := range files {
        path := filepath.Join(migrationsDir, f.Name())
        sqlBytes, err := os.ReadFile(path)
        if err != nil {
            return err
        }
        log.Printf("Applying migration: %s", f.Name())
        if _, err := db.Exec(string(sqlBytes)); err != nil {
            return err
        }
    }
    log.Printf("Migrations applied successfully")
    return nil
}