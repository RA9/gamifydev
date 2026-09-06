package store

import (
	"context"
	"path/filepath"
	"testing"
)

func TestCurriculumMigrationArchivesPythonCheckpointsAndResetsFoundationsSchedule(t *testing.T) {
	ctx := context.Background()
	st, err := Open(filepath.Join(t.TempDir(), "upgrade.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	mustExec(t, st, `CREATE TABLE schema_migrations (name TEXT PRIMARY KEY, applied_at TEXT NOT NULL DEFAULT (datetime('now')))`)
	mustExec(t, st, `INSERT INTO schema_migrations (name) VALUES ('0029_curriculum_model.sql')`)
	if err := st.Migrate(ctx); err != nil {
		t.Fatalf("migrate through phase 2: %v", err)
	}
	mustExec(t, st, `INSERT INTO paths (id, slug, title) VALUES (90, 'cs-foundations', 'Foundations')`)
	mustExec(t, st, `INSERT INTO cohorts (id, path_id, name, tz_band, starts_on, scheduled_at)
		VALUES (90, 90, 'Legacy Foundations', 'europe_africa', '2026-01-05', datetime('now'))`)
	mustExec(t, st, `INSERT INTO cohort_schedule
		(cohort_id, day_index, due_on, kind, sort) VALUES (90, 0, '2026-01-05', 'lesson', 0)`)
	mustExec(t, st, `INSERT INTO assignments
		(slug, title, language, published, required) VALUES ('ds-hash-map', 'Old Python Map', 'python', 1, 1)`)

	mustExec(t, st, `DELETE FROM schema_migrations WHERE name = '0029_curriculum_model.sql'`)
	if err := st.Migrate(ctx); err != nil {
		t.Fatalf("apply curriculum migration: %v", err)
	}
	var scheduled int
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM cohort_schedule WHERE cohort_id = 90`).Scan(&scheduled); err != nil {
		t.Fatalf("count legacy schedule: %v", err)
	}
	if scheduled != 0 {
		t.Fatalf("legacy Foundations schedule retained %d item(s)", scheduled)
	}
	var scheduledAt any
	if err := st.db.QueryRow(`SELECT scheduled_at FROM cohorts WHERE id = 90`).Scan(&scheduledAt); err != nil {
		t.Fatalf("read cohort: %v", err)
	}
	if scheduledAt != nil {
		t.Fatalf("legacy Foundations cohort still marked scheduled: %v", scheduledAt)
	}
	var published, required int
	if err := st.db.QueryRow(`SELECT published, required FROM assignments WHERE slug='ds-hash-map'`).Scan(&published, &required); err != nil {
		t.Fatalf("read archived assignment: %v", err)
	}
	if published != 0 || required != 0 {
		t.Fatalf("legacy Python checkpoint remained active: published=%d required=%d", published, required)
	}
}
