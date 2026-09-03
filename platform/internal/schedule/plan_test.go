package schedule

import (
	"testing"
	"time"
)

func lessons(ids ...int64) []Lesson {
	var out []Lesson
	for _, id := range ids {
		out = append(out, Lesson{ID: id})
	}
	return out
}

// monday is a known Monday, so weekday arithmetic in tests is readable.
var monday = time.Date(2026, 9, 7, 0, 0, 0, 0, time.UTC)

func TestPlanFillsTheDailyBudget(t *testing.T) {
	items := Plan([]Course{{ID: 1, Lessons: lessons(1, 2, 3, 4, 5)}}, monday)
	if len(items) != 5 {
		t.Fatalf("got %d items, want 5", len(items))
	}
	perDay := map[int]int{}
	for _, it := range items {
		perDay[it.DayIndex]++
	}
	for day, n := range perDay {
		if n > LessonsPerDay {
			t.Fatalf("day %d has %d lessons, over the budget of %d", day, n, LessonsPerDay)
		}
	}
	// 5 lessons at 2/day is 3 days.
	if got := Days(items); got != 3 {
		t.Fatalf("plan spans %d days, want 3", got)
	}
}

func TestPlanSkipsWeekends(t *testing.T) {
	// Enough lessons to run past a Friday.
	items := Plan([]Course{{ID: 1, Lessons: lessons(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)}}, monday)
	for _, it := range items {
		switch it.DueOn.Weekday() {
		case time.Saturday, time.Sunday:
			t.Fatalf("item scheduled on a %s (%s)", it.DueOn.Weekday(), it.DueOn.Format("2006-01-02"))
		}
	}
	// Day 5 is the second Monday, not Saturday.
	var day5 time.Time
	for _, it := range items {
		if it.DayIndex == 5 {
			day5 = it.DueOn
		}
	}
	if day5.Weekday() != time.Monday {
		t.Fatalf("working day 5 fell on %s, want Monday", day5.Weekday())
	}
	if want := monday.AddDate(0, 0, 7); !day5.Equal(want) {
		t.Fatalf("working day 5 = %s, want %s", day5.Format("2006-01-02"), want.Format("2006-01-02"))
	}
}

func TestPlanStartingOnAWeekendBeginsMonday(t *testing.T) {
	saturday := time.Date(2026, 9, 5, 0, 0, 0, 0, time.UTC)
	if saturday.Weekday() != time.Saturday {
		t.Fatalf("test fixture is a %s", saturday.Weekday())
	}
	items := Plan([]Course{{ID: 1, Lessons: lessons(1)}}, saturday)
	if got := items[0].DueOn.Weekday(); got != time.Monday {
		t.Fatalf("a cohort starting Saturday has work due on %s, want Monday", got)
	}
}

func TestCoursesDoNotShareADay(t *testing.T) {
	// An odd lesson count leaves a half-used day; the next course must still
	// start fresh so "today" is describable in one line.
	items := Plan([]Course{
		{ID: 1, Lessons: lessons(1)},
		{ID: 2, Lessons: lessons(2)},
	}, monday)
	byDay := map[int]map[int64]bool{}
	for _, it := range items {
		if byDay[it.DayIndex] == nil {
			byDay[it.DayIndex] = map[int64]bool{}
		}
		byDay[it.DayIndex][it.CourseID] = true
	}
	for day, courses := range byDay {
		if len(courses) > 1 {
			t.Fatalf("day %d mixes %d courses", day, len(courses))
		}
	}
}

func TestCheckpointLandsWithItsCourse(t *testing.T) {
	items := Plan([]Course{
		{ID: 1, Lessons: lessons(1, 2, 3), CheckpointID: 100},
		{ID: 2, Lessons: lessons(4)},
	}, monday)

	var checkpoint, lastLesson Item
	for _, it := range items {
		if it.Kind == KindCheckpoint {
			checkpoint = it
		}
		if it.Kind == KindLesson && it.CourseID == 1 && it.DayIndex >= lastLesson.DayIndex {
			lastLesson = it
		}
	}
	if checkpoint.AssignmentID != 100 {
		t.Fatal("checkpoint not scheduled")
	}
	// Due the same day the course's material finishes — not drifting into the
	// next course's week.
	if checkpoint.DayIndex != lastLesson.DayIndex {
		t.Fatalf("checkpoint on day %d, course's last lesson on day %d",
			checkpoint.DayIndex, lastLesson.DayIndex)
	}
	if checkpoint.Sort <= lastLesson.Sort {
		t.Fatalf("checkpoint sorts before the lesson it follows (%d vs %d)",
			checkpoint.Sort, lastLesson.Sort)
	}
	// And the next course starts after it.
	for _, it := range items {
		if it.CourseID == 2 && it.DayIndex <= checkpoint.DayIndex {
			t.Fatalf("course 2 starts on day %d, not after the checkpoint on day %d",
				it.DayIndex, checkpoint.DayIndex)
		}
	}
}

func TestSprintsGroupWorkingDays(t *testing.T) {
	// Ten working days at 2/day = 20 lessons, which should be two sprints.
	var ids []int64
	for i := int64(1); i <= 20; i++ {
		ids = append(ids, i)
	}
	items := Plan([]Course{{ID: 1, Lessons: lessons(ids...)}}, monday)
	if got := Days(items); got != 10 {
		t.Fatalf("plan spans %d days, want 10", got)
	}
	if got := Sprints(items); got != 2 {
		t.Fatalf("plan spans %d sprints, want 2", got)
	}
	for _, it := range items {
		want := it.DayIndex/DaysPerSprint + 1
		if it.Sprint != want {
			t.Fatalf("day %d in sprint %d, want %d", it.DayIndex, it.Sprint, want)
		}
	}
}

func TestEmptyInputs(t *testing.T) {
	if items := Plan(nil, monday); len(items) != 0 {
		t.Fatalf("nil courses produced %d items", len(items))
	}
	// A course with neither lessons nor a checkpoint contributes nothing and
	// must not consume a day.
	items := Plan([]Course{{ID: 1}, {ID: 2, Lessons: lessons(1)}}, monday)
	if len(items) != 1 {
		t.Fatalf("got %d items, want 1", len(items))
	}
	if items[0].DayIndex != 0 {
		t.Fatalf("an empty course consumed a day (item on day %d)", items[0].DayIndex)
	}
	if got := Days(nil); got != 0 {
		t.Fatalf("Days(nil) = %d, want 0", got)
	}
	if got := Sprints(nil); got != 0 {
		t.Fatalf("Sprints(nil) = %d, want 0", got)
	}
}

func TestPlanIsDeterministic(t *testing.T) {
	courses := []Course{
		{ID: 1, Lessons: lessons(1, 2, 3), CheckpointID: 10},
		{ID: 2, Lessons: lessons(4, 5)},
	}
	first := Plan(courses, monday)
	for i := 0; i < 5; i++ {
		got := Plan(courses, monday)
		if len(got) != len(first) {
			t.Fatalf("plan length varies: %d then %d", len(first), len(got))
		}
		for j := range got {
			if got[j] != first[j] {
				t.Fatalf("item %d differs between runs:\n %+v\n %+v", j, first[j], got[j])
			}
		}
	}
}
