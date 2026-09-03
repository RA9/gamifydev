// Package schedule turns a path's contents into a cohort's dated plan.
//
// Pure, like internal/placement and internal/cohort: given the ordered courses,
// their lessons, and their checkpoints, produce the day-indexed items. No
// database, no clock beyond the start date it is handed — which is what makes
// the pacing rules cheap to change and cheap to test.
//
// This is the phase that ends self-pacing (PRD §07). The cadence numbers below
// are the whole policy, so they live together and are asserted on directly.
package schedule

import "time"

const (
	// LessonsPerDay is the daily budget. Two substantial lessons is roughly an
	// evening's work, which is the realistic ceiling for people fitting this
	// around a job or school.
	LessonsPerDay = 2

	// DaysPerSprint is how many working days make a "week" for display.
	DaysPerSprint = 5
)

// Kinds of scheduled item.
const (
	KindLesson     = "lesson"
	KindCheckpoint = "checkpoint"
)

// Lesson is the minimum a lesson contributes to planning.
type Lesson struct {
	ID    int64
	Title string
}

// Course is one path step: its lessons in order, plus an optional checkpoint
// assignment that gates progression to the next course.
type Course struct {
	ID           int64
	Title        string
	Slug         string
	Lessons      []Lesson
	CheckpointID int64 // 0 when the course has no required assignment
}

// Item is one scheduled piece of work.
type Item struct {
	DayIndex     int
	Sprint       int
	DueOn        time.Time
	Kind         string
	CourseID     int64
	LessonID     int64
	AssignmentID int64
	Sort         int
}

// Plan expands a path into dated items starting on `start`.
//
// Lessons fill working days at LessonsPerDay each, in path order. A course's
// checkpoint lands on the last day that course occupies, so it is due when the
// material is done rather than drifting into the next course's week.
//
// Weekends are skipped. A plan that demands seven days a week is one nobody
// sustains, and — since attendance is measured against scheduled days — it would
// make the absence rules far harsher than intended.
func Plan(courses []Course, start time.Time) []Item {
	start = start.UTC().Truncate(24 * time.Hour)
	var items []Item

	day := 0       // working-day index from the start
	usedToday := 0 // lessons placed on the current day
	sortInDay := 0 // ordering within a day

	// advance moves to the next working day.
	advance := func() {
		day++
		usedToday = 0
		sortInDay = 0
	}

	for _, c := range courses {
		if len(c.Lessons) == 0 && c.CheckpointID == 0 {
			continue
		}
		// A new course starts on a fresh day so a day's work never spans two
		// courses — it keeps "today" describable in one line.
		if usedToday > 0 {
			advance()
		}
		for _, l := range c.Lessons {
			if usedToday >= LessonsPerDay {
				advance()
			}
			items = append(items, Item{
				DayIndex: day, Sprint: sprintOf(day), DueOn: workingDay(start, day),
				Kind: KindLesson, CourseID: c.ID, LessonID: l.ID, Sort: sortInDay,
			})
			usedToday++
			sortInDay++
		}
		if c.CheckpointID != 0 {
			// Due on the same day the course's last lesson lands, after it.
			if len(c.Lessons) == 0 {
				usedToday = 0
			}
			items = append(items, Item{
				DayIndex: day, Sprint: sprintOf(day), DueOn: workingDay(start, day),
				Kind: KindCheckpoint, CourseID: c.ID, AssignmentID: c.CheckpointID,
				Sort: sortInDay,
			})
			sortInDay++
			// Force the next course onto a new day.
			usedToday = LessonsPerDay
		}
	}
	return items
}

// sprintOf is the 1-based sprint (week) a working-day index falls in.
func sprintOf(dayIndex int) int { return dayIndex/DaysPerSprint + 1 }

// workingDay returns the calendar date `n` working days after `start`,
// skipping weekends. If `start` itself is a weekend, the plan begins on the
// following Monday.
func workingDay(start time.Time, n int) time.Time {
	d := start
	for isWeekend(d) {
		d = d.AddDate(0, 0, 1)
	}
	for i := 0; i < n; i++ {
		d = d.AddDate(0, 0, 1)
		for isWeekend(d) {
			d = d.AddDate(0, 0, 1)
		}
	}
	return d
}

func isWeekend(t time.Time) bool {
	switch t.Weekday() {
	case time.Saturday, time.Sunday:
		return true
	}
	return false
}

// Days reports how many working days a plan spans.
func Days(items []Item) int {
	max := -1
	for _, it := range items {
		if it.DayIndex > max {
			max = it.DayIndex
		}
	}
	return max + 1
}

// Sprints reports how many sprints a plan spans.
func Sprints(items []Item) int {
	if len(items) == 0 {
		return 0
	}
	return sprintOf(Days(items) - 1)
}
