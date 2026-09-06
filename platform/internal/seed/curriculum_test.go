package seed

import (
	"context"
	"database/sql"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/RA9/gamifydev/platform/internal/store"
)

func seededCurriculumDB(t *testing.T) *sql.DB {
	t.Helper()
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "curriculum.db")
	st, err := store.Open(path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	if err := st.Migrate(ctx); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if _, err := Run(ctx, st); err != nil {
		t.Fatalf("seed: %v", err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open verification database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func TestFoundationsWorkloadMatchesProgramRatio(t *testing.T) {
	db := seededCurriculumDB(t)
	rows, err := db.Query(`
		SELECT mode, SUM(workload) FROM (
		  SELECT l.mode, l.workload_minutes AS workload
		  FROM paths p
		  JOIN path_courses pc ON pc.path_id = p.id
		  JOIN lessons l ON l.course_id = pc.course_id
		  WHERE p.slug = 'cs-foundations'
		  UNION ALL
		  SELECT a.mode, a.workload_minutes
		  FROM paths p
		  JOIN path_courses pc ON pc.path_id = p.id
		  JOIN assignments a ON a.course_id = pc.course_id
		  WHERE p.slug = 'cs-foundations' AND a.published = 1
		  UNION ALL
		  SELECT pr.mode, pr.workload_minutes
		  FROM paths p
		  JOIN path_courses pc ON pc.path_id = p.id
		  JOIN course_problems cp ON cp.course_id = pc.course_id
		  JOIN problems pr ON pr.id = cp.problem_id
		  WHERE p.slug = 'cs-foundations' AND pr.published = 1 AND cp.required = 1
		) GROUP BY mode`)
	if err != nil {
		t.Fatalf("load workload: %v", err)
	}
	defer rows.Close()
	byMode := map[string]int{}
	total := 0
	for rows.Next() {
		var mode string
		var workload int
		if err := rows.Scan(&mode, &workload); err != nil {
			t.Fatalf("scan workload: %v", err)
		}
		byMode[mode] = workload
		total += workload
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("workload rows: %v", err)
	}
	if total == 0 {
		t.Fatal("Foundations has no workload")
	}
	want := map[string]int{
		store.ModeFun:         15,
		store.ModeTheoretical: 30,
		store.ModePractical:   55,
	}
	for mode, percent := range want {
		if byMode[mode]*100 != total*percent {
			t.Errorf("%s workload = %d/%d (%d%% target)", mode, byMode[mode], total, percent)
		}
	}
	if total != 2000 {
		t.Errorf("Foundations workload = %d minutes, want 2000", total)
	}
}

func TestFoundationsUsesCWithOnlyDocumentedLinuxExceptions(t *testing.T) {
	db := seededCurriculumDB(t)
	var javaCount int
	if err := db.QueryRow(`
		SELECT COUNT(*) FROM paths p
		JOIN path_courses pc ON pc.path_id = p.id
		JOIN courses c ON c.id = pc.course_id
		WHERE p.slug = 'cs-foundations' AND c.slug = 'java'`).Scan(&javaCount); err != nil {
		t.Fatalf("check Java membership: %v", err)
	}
	if javaCount != 0 {
		t.Fatal("Java is still part of the Foundations path")
	}

	courseRows, err := db.Query(`
		SELECT c.slug, c.primary_language, c.language_policy, c.language_exception
		FROM paths p JOIN path_courses pc ON pc.path_id=p.id JOIN courses c ON c.id=pc.course_id
		WHERE p.slug='cs-foundations'`)
	if err != nil {
		t.Fatalf("load course language policies: %v", err)
	}
	for courseRows.Next() {
		var course, language, policy, exception string
		if err := courseRows.Scan(&course, &language, &policy, &exception); err != nil {
			courseRows.Close()
			t.Fatalf("scan course language policy: %v", err)
		}
		if course == "linux" {
			if language != "shell" || policy != store.LanguagePolicyCWithException || strings.TrimSpace(exception) == "" {
				t.Errorf("Linux language policy = %q/%q/%q", language, policy, exception)
			}
		} else if language != "c" || policy != store.LanguagePolicyCOnly || exception != "" {
			t.Errorf("%s language policy = %q/%q/%q, want c/c_only/no exception", course, language, policy, exception)
		}
	}
	if err := courseRows.Err(); err != nil {
		courseRows.Close()
		t.Fatalf("course language policies: %v", err)
	}
	courseRows.Close()

	queries := []struct {
		name, query string
	}{
		{"lessons", `SELECT c.slug, l.slug, l.language FROM paths p JOIN path_courses pc ON pc.path_id=p.id JOIN courses c ON c.id=pc.course_id JOIN lessons l ON l.course_id=c.id WHERE p.slug='cs-foundations' AND ((c.slug='linux' AND l.language!='shell') OR (c.slug!='linux' AND l.language!='c'))`},
		{"checkpoints", `SELECT c.slug, a.slug, a.language FROM paths p JOIN path_courses pc ON pc.path_id=p.id JOIN courses c ON c.id=pc.course_id JOIN assignments a ON a.course_id=c.id WHERE p.slug='cs-foundations' AND a.required=1 AND ((c.slug='linux' AND a.language!='shell') OR (c.slug!='linux' AND a.language!='c'))`},
		{"problems", `SELECT c.slug, pr.slug, pr.language FROM paths p JOIN path_courses pc ON pc.path_id=p.id JOIN courses c ON c.id=pc.course_id JOIN course_problems cp ON cp.course_id=c.id JOIN problems pr ON pr.id=cp.problem_id WHERE p.slug='cs-foundations' AND pr.language!='c'`},
		{"steps", `SELECT c.slug, l.slug, s.lang FROM paths p JOIN path_courses pc ON pc.path_id=p.id JOIN courses c ON c.id=pc.course_id JOIN lessons l ON l.course_id=c.id JOIN lesson_steps s ON s.lesson_id=l.id WHERE p.slug='cs-foundations' AND ((c.slug='linux' AND s.lang!='shell') OR (c.slug!='linux' AND s.lang!='c'))`},
	}
	for _, check := range queries {
		rows, err := db.Query(check.query)
		if err != nil {
			t.Fatalf("check %s languages: %v", check.name, err)
		}
		for rows.Next() {
			var course, item, language string
			if err := rows.Scan(&course, &item, &language); err != nil {
				rows.Close()
				t.Fatalf("scan %s violation: %v", check.name, err)
			}
			t.Errorf("%s %s/%s uses %q", check.name, course, item, language)
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			t.Fatalf("%s language rows: %v", check.name, err)
		}
		rows.Close()
	}

	var exception string
	if err := db.QueryRow(`SELECT language_exception FROM courses WHERE slug='linux'`).Scan(&exception); err != nil {
		t.Fatalf("load Linux exception: %v", err)
	}
	if strings.TrimSpace(exception) == "" {
		t.Fatal("Linux uses shell but has no documented language exception")
	}
}

func TestFoundationsNotesContainNoPythonOrJavaCodeFences(t *testing.T) {
	for course := range foundationCourses {
		dir := "data/notes/" + course
		err := fs.WalkDir(contentFS, dir, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			body, err := contentFS.ReadFile(path)
			if err != nil {
				return err
			}
			lower := strings.ToLower(string(body))
			for _, forbidden := range []string{"```python", "```py\n", "```java", "```javascript"} {
				if strings.Contains(lower, forbidden) {
					t.Errorf("%s contains forbidden Foundations fence %q", path, forbidden)
				}
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", dir, err)
		}
	}
}

func TestEveryFoundationsCourseIntegratesALabProblemAndCheckpoint(t *testing.T) {
	db := seededCurriculumDB(t)
	rows, err := db.Query(`
		SELECT c.slug,
		       (SELECT COUNT(*) FROM lessons l WHERE l.course_id=c.id AND l.kind='lab'),
		       (SELECT COUNT(*) FROM course_problems cp WHERE cp.course_id=c.id),
		       (SELECT COUNT(*) FROM assignments a WHERE a.course_id=c.id AND a.required=1)
		FROM paths p JOIN path_courses pc ON pc.path_id=p.id JOIN courses c ON c.id=pc.course_id
		WHERE p.slug='cs-foundations'`)
	if err != nil {
		t.Fatalf("load Foundations integration: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var course string
		var labs, problems, checkpoints int
		if err := rows.Scan(&course, &labs, &problems, &checkpoints); err != nil {
			t.Fatalf("scan course integration: %v", err)
		}
		if labs == 0 || problems == 0 || checkpoints == 0 {
			t.Errorf("%s has labs=%d problems=%d checkpoints=%d", course, labs, problems, checkpoints)
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("course integration rows: %v", err)
	}
}
