package store

import (
	"context"
	"errors"
	"testing"
)

func TestCourseLanguagePolicyValidation(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	for _, course := range []Course{
		{Slug: "bad-policy", Title: "Bad", LanguagePolicy: "anything"},
		{Slug: "missing-reason", Title: "Missing", PrimaryLanguage: "c", LanguagePolicy: LanguagePolicyCWithException},
	} {
		if _, err := st.UpsertCourse(ctx, course); !errors.Is(err, ErrInvalidCourseLanguagePolicy) {
			t.Fatalf("policy %q err = %v, want ErrInvalidCourseLanguagePolicy", course.LanguagePolicy, err)
		}
	}
	id, err := st.UpsertCourse(ctx, Course{
		Slug: "justified", Title: "Justified", PrimaryLanguage: "shell",
		LanguagePolicy: LanguagePolicyCWithException, LanguageException: "Linux operations require shell.",
	})
	if err != nil {
		t.Fatalf("valid exception policy: %v", err)
	}
	course, err := st.GetCourseByID(ctx, id)
	if err != nil || course.LanguageException == "" {
		t.Fatalf("policy did not round-trip: %+v, %v", course, err)
	}
}

func TestReplaceLessonStepsPreservesProgressForStablePositions(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	uid := mkUser(t, st, "steps@example.com")
	mustExec(t, st, `INSERT INTO courses (id, slug, title) VALUES (50, 'lab-course', 'Lab')`)
	mustExec(t, st, `INSERT INTO lessons (id, course_id, slug, title) VALUES (50, 50, 'lab', 'Lab')`)
	first := []Step{
		{Lang: "c", Instruction: "one", Starter: "a", Checks: "[]"},
		{Lang: "c", Instruction: "two", Starter: "b", Checks: "[]"},
	}
	if err := st.ReplaceLessonSteps(ctx, 50, first); err != nil {
		t.Fatalf("first replace: %v", err)
	}
	steps, err := st.ListSteps(ctx, 50)
	if err != nil || len(steps) != 2 {
		t.Fatalf("first steps = %v, %v", steps, err)
	}
	firstID := steps[0].ID
	if err := st.MarkStepComplete(ctx, uid, firstID); err != nil {
		t.Fatalf("mark complete: %v", err)
	}

	updated := []Step{
		{Lang: "c", Instruction: "one updated", Starter: "updated", Checks: "[]"},
		{Lang: "shell", Instruction: "two updated", Starter: "echo ok", Checks: "[]"},
	}
	if err := st.ReplaceLessonSteps(ctx, 50, updated); err != nil {
		t.Fatalf("second replace: %v", err)
	}
	steps, err = st.ListSteps(ctx, 50)
	if err != nil || len(steps) != 2 {
		t.Fatalf("updated steps = %v, %v", steps, err)
	}
	if steps[0].ID != firstID || steps[0].Instruction != "one updated" || steps[0].Lang != "c" || steps[1].Lang != "shell" {
		t.Fatalf("stable step rows were not updated in place: %+v", steps)
	}
	completed, err := st.CompletedStepIDs(ctx, uid, 50)
	if err != nil || !completed[firstID] {
		t.Fatalf("reseeding lost progress: %v, %v", completed, err)
	}
}
