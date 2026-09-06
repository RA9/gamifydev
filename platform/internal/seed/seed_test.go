package seed

import (
	"context"
	"database/sql"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/RA9/gamifydev/platform/internal/store"
)

// TestRunIsIdempotent guards the startup contract: the server runs the complete
// content seed on every boot, so a second run must update the existing rows
// rather than append another copy of any seeded record or child collection.
func TestRunIsIdempotent(t *testing.T) {
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "seed.db")

	st, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })

	if err := st.Migrate(ctx); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if _, err := Run(ctx, st); err != nil {
		t.Fatalf("first seed: %v", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open count connection: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	before := seededTableCounts(t, db)
	if _, err := Run(ctx, st); err != nil {
		t.Fatalf("second seed: %v", err)
	}
	after := seededTableCounts(t, db)

	if !reflect.DeepEqual(after, before) {
		t.Fatalf("second seed changed row counts\nbefore: %#v\nafter:  %#v", before, after)
	}
}

func seededTableCounts(t *testing.T, db *sql.DB) map[string]int {
	t.Helper()
	tables := []string{
		"courses",
		"lessons",
		"assignments",
		"assignment_checks",
		"assignment_files",
		"lesson_steps",
		"paths",
		"path_courses",
		"assessments",
		"assessment_items",
		"problems",
		"problem_starters",
		"problem_tests",
		"course_problems",
	}
	counts := make(map[string]int, len(tables))
	for _, table := range tables {
		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
			t.Fatalf("count %s: %v", table, err)
		}
		counts[table] = count
	}
	return counts
}
