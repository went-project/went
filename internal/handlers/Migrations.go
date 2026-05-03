package handlers

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

// migrationRunner wraps a sql.DB with dialect info for portable queries.
type migrationRunner struct {
	db      *sql.DB
	dialect string
}

func newMigrationRunner() (*migrationRunner, error) {
	_ = godotenv.Load(".env.local")
	_ = godotenv.Load(".env")

	dialect := envOr("DB_DIALECT", "sqlite")
	host := envOr("DB_HOST", "localhost")
	port := envOr("DB_PORT", "5432")
	user := envOr("DB_USER", "postgres")
	pass := envOr("DB_PASSWORD", "password")
	dbname := envOr("DB_NAME", "")

	if dialect == "sqlite" && (dbname == "" || dbname == "database") {
		dbname = envOr("DB_STORAGE", "./database.sqlite")
	}

	var (
		db  *sql.DB
		err error
	)

	switch dialect {
	case "sqlite":
		db, err = sql.Open("sqlite", dbname)
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			user, pass, host, port, dbname)
		db, err = sql.Open("mysql", dsn)
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			host, port, user, dbname, pass)
		db, err = sql.Open("postgres", dsn)
	default:
		return nil, fmt.Errorf("unsupported dialect: %s", dialect)
	}

	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	return &migrationRunner{db: db, dialect: dialect}, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func (r *migrationRunner) close() {
	r.db.Close()
}

func (r *migrationRunner) ensureTable() error {
	_, err := r.db.Exec(`CREATE TABLE IF NOT EXISTS wentmigrations (
		migration VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	return err
}

func (r *migrationRunner) appliedMap() (map[string]bool, error) {
	rows, err := r.db.Query("SELECT migration FROM wentmigrations ORDER BY migration ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]bool)
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		result[s] = true
	}
	return result, rows.Err()
}

func (r *migrationRunner) appliedOrdered() ([]string, error) {
	rows, err := r.db.Query("SELECT migration FROM wentmigrations ORDER BY migration ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, rows.Err()
}

func (r *migrationRunner) record(name string) error {
	var err error
	if r.dialect == "postgres" {
		_, err = r.db.Exec("INSERT INTO wentmigrations (migration) VALUES ($1)", name)
	} else {
		_, err = r.db.Exec("INSERT INTO wentmigrations (migration) VALUES (?)", name)
	}
	return err
}

func (r *migrationRunner) unrecord(name string) error {
	var err error
	if r.dialect == "postgres" {
		_, err = r.db.Exec("DELETE FROM wentmigrations WHERE migration = $1", name)
	} else {
		_, err = r.db.Exec("DELETE FROM wentmigrations WHERE migration = ?", name)
	}
	return err
}

// --- file helpers ---

func upFiles() ([]string, error) {
	files, err := filepath.Glob(filepath.Join("database", "migrations", "*.up.sql"))
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

func baseName(upFile string) string {
	return strings.TrimSuffix(filepath.Base(upFile), ".up.sql")
}

func downFilePath(name string) string {
	return filepath.Join("database", "migrations", name+".down.sql")
}

// --- public commands ---

// RunMigrations runs all pending up migrations in order.
func RunMigrations() error {
	runner, err := newMigrationRunner()
	if err != nil {
		return err
	}
	defer runner.close()

	if err := runner.ensureTable(); err != nil {
		return err
	}

	applied, err := runner.appliedMap()
	if err != nil {
		return err
	}

	files, err := upFiles()
	if err != nil {
		return err
	}

	ran := 0
	for _, f := range files {
		name := baseName(f)
		if applied[name] {
			continue
		}

		content, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("read %s: %w", f, err)
		}
		if _, err := runner.db.Exec(string(content)); err != nil {
			return fmt.Errorf("run %s: %w", name, err)
		}
		if err := runner.record(name); err != nil {
			return fmt.Errorf("record %s: %w", name, err)
		}
		fmt.Printf("  ✓ Migrated:   %s\n", name)
		ran++
	}

	if ran == 0 {
		fmt.Println("Nothing to migrate.")
	}
	return nil
}

// RollbackMigrations rolls back the last `step` applied migrations.
func RollbackMigrations(step int) error {
	runner, err := newMigrationRunner()
	if err != nil {
		return err
	}
	defer runner.close()

	if err := runner.ensureTable(); err != nil {
		return err
	}

	ordered, err := runner.appliedOrdered()
	if err != nil {
		return err
	}
	if len(ordered) == 0 {
		fmt.Println("Nothing to rollback.")
		return nil
	}

	if step > len(ordered) {
		step = len(ordered)
	}

	// take last `step` entries and reverse for rollback order
	toRollback := make([]string, step)
	copy(toRollback, ordered[len(ordered)-step:])
	for i, j := 0, len(toRollback)-1; i < j; i, j = i+1, j-1 {
		toRollback[i], toRollback[j] = toRollback[j], toRollback[i]
	}

	for _, name := range toRollback {
		df := downFilePath(name)
		content, err := os.ReadFile(df)
		if err != nil {
			return fmt.Errorf("read %s: %w", df, err)
		}
		if _, err := runner.db.Exec(string(content)); err != nil {
			return fmt.Errorf("rollback %s: %w", name, err)
		}
		if err := runner.unrecord(name); err != nil {
			return fmt.Errorf("unrecord %s: %w", name, err)
		}
		fmt.Printf("  ✓ Rolled back: %s\n", name)
	}
	return nil
}

// FreshMigrations drops all migrated tables then re-runs every migration.
func FreshMigrations() error {
	runner, err := newMigrationRunner()
	if err != nil {
		return err
	}
	defer runner.close()

	// collect all up files to get their down counterparts
	files, err := upFiles()
	if err != nil {
		return err
	}

	// run down files in reverse order; ignore errors (table may not exist)
	for i := len(files) - 1; i >= 0; i-- {
		name := baseName(files[i])
		df := downFilePath(name)
		content, err := os.ReadFile(df)
		if err != nil {
			continue
		}
		_, _ = runner.db.Exec(string(content))
		fmt.Printf("  ↓ Dropped:    %s\n", name)
	}

	_, _ = runner.db.Exec("DROP TABLE IF EXISTS wentmigrations")

	// recreate tracking table and run all up files
	if err := runner.ensureTable(); err != nil {
		return err
	}

	for _, f := range files {
		name := baseName(f)
		content, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("read %s: %w", f, err)
		}
		if _, err := runner.db.Exec(string(content)); err != nil {
			return fmt.Errorf("run %s: %w", name, err)
		}
		if err := runner.record(name); err != nil {
			return fmt.Errorf("record %s: %w", name, err)
		}
		fmt.Printf("  ✓ Migrated:   %s\n", name)
	}
	return nil
}
