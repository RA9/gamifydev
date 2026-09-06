package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/content"
	"github.com/RA9/gamifydev/platform/internal/store"
)

// --- Learner / public course delivery --------------------------------------

func (s *Server) handleCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := s.st.ListCourses(r.Context(), false)
	if err != nil {
		http.Error(w, "could not load courses", http.StatusInternalServerError)
		return
	}
	s.render(w, r, "courses.html", ViewData{Title: "Courses", Data: map[string]any{"courses": courses}})
}

func (s *Server) handleCourse(w http.ResponseWriter, r *http.Request) {
	course, err := s.st.GetCourseBySlug(r.Context(), r.PathValue("course"))
	if err != nil {
		s.notFound(w, r)
		return
	}
	lessons, _ := s.st.ListLessons(r.Context(), course.ID)
	assignments, _ := s.st.ListAssignmentsByCourse(r.Context(), course.ID, false)
	problems, _ := s.st.ListProblemsByCourse(r.Context(), course.ID, s.solver(r))
	u := auth.CurrentUser(r.Context())
	locked := u == nil
	passed := map[int64]bool{}
	var uid int64
	if u != nil {
		uid = u.ID
		passed, _ = s.st.AssignmentPassState(r.Context(), u.ID, course.ID)
	}
	stepProgress, _ := s.st.StepProgressByLesson(r.Context(), uid, course.ID)
	s.render(w, r, "course.html", ViewData{
		Title: course.Title,
		Data: map[string]any{
			"course":       course,
			"sections":     groupLessons(lessons),
			"lessonCount":  len(lessons),
			"assignments":  assignments,
			"problems":     problems,
			"locked":       locked,
			"passed":       passed,
			"stepProgress": stepProgress,
		},
	})
}

// lessonSection is a named group of lessons within a course.
type lessonSection struct {
	Name    string
	Lessons []store.Lesson
}

// groupLessons groups lessons into their sections, preserving lesson order (so a
// section appears at the position of its first lesson).
func groupLessons(lessons []store.Lesson) []lessonSection {
	var out []lessonSection
	idx := map[string]int{}
	for _, l := range lessons {
		i, ok := idx[l.Section]
		if !ok {
			out = append(out, lessonSection{Name: l.Section})
			i = len(out) - 1
			idx[l.Section] = i
		}
		out[i].Lessons = append(out[i].Lessons, l)
	}
	return out
}

func (s *Server) handleLesson(w http.ResponseWriter, r *http.Request) {
	course, err := s.st.GetCourseBySlug(r.Context(), r.PathValue("course"))
	if err != nil {
		s.notFound(w, r)
		return
	}
	lesson, err := s.st.GetLesson(r.Context(), course.ID, r.PathValue("lesson"))
	if err != nil {
		s.notFound(w, r)
		return
	}
	// Determine prev/next within the course for sequential navigation.
	lessons, _ := s.st.ListLessons(r.Context(), course.ID)
	var prev, next *store.Lesson
	for i := range lessons {
		if lessons[i].ID == lesson.ID {
			if i > 0 {
				prev = &lessons[i-1]
			}
			if i < len(lessons)-1 {
				next = &lessons[i+1]
			}
			break
		}
	}
	// Guests can see the lesson exists, but content is gated behind sign-in.
	u := auth.CurrentUser(r.Context())
	locked := u == nil
	if u != nil && writeWorkAccessError(w, s.st.RequireLessonAvailable(r.Context(), u.ID, lesson.ID)) {
		return
	}
	steps, _ := s.st.ListSteps(r.Context(), lesson.ID)

	data := map[string]any{
		"course": course,
		"lesson": lesson,
		"prev":   prev,
		"next":   next,
		"locked": locked,
	}

	switch {
	case locked:
		data["hasSteps"] = len(steps) > 0
	case len(steps) > 0:
		s.renderStepLab(w, r, course, lesson, steps, next)
		return
	default:
		data["body"] = content.Render(lesson.Body)
		// A reading lesson has no steps to finish, so it needs an explicit
		// "done" — otherwise it can never be cleared off the schedule.
		data["done"], _ = s.st.LessonComplete(r.Context(), u.ID, lesson.ID)
		data["completeURL"] = "/lessons/" + strconv.FormatInt(lesson.ID, 10) + "/complete"
	}
	s.render(w, r, "lesson.html", ViewData{Title: lesson.Title, Data: data})
}

type stepCheck struct {
	Text string `json:"text"`
	Test string `json:"test"`
}

// renderStepLab renders the interactive, step-by-step view of a lesson for a
// signed-in learner.
func (s *Server) renderStepLab(w http.ResponseWriter, r *http.Request, course *store.Course, lesson *store.Lesson, steps []store.Step, next *store.Lesson) {
	u := auth.CurrentUser(r.Context())
	completed, _ := s.st.CompletedStepIDs(r.Context(), u.ID, lesson.ID)

	// Current step: explicit ?step=N (1-based), else first incomplete, else last.
	cur := -1
	if n, err := strconv.Atoi(r.URL.Query().Get("step")); err == nil && n >= 1 && n <= len(steps) {
		cur = n - 1
	}
	if cur == -1 {
		for i, st := range steps {
			if !completed[st.ID] {
				cur = i
				break
			}
		}
		if cur == -1 {
			cur = len(steps) - 1
		}
	}
	step := steps[cur]

	var checks []stepCheck
	_ = json.Unmarshal([]byte(step.Checks), &checks)

	nextURL := ""
	isLast := cur == len(steps)-1
	if !isLast {
		nextURL = "?step=" + strconv.Itoa(cur+2)
	} else if next != nil {
		nextURL = "/courses/" + course.Slug + "/" + next.Slug
	} else {
		nextURL = "/courses/" + course.Slug
	}
	prevURL := ""
	if cur > 0 {
		prevURL = "?step=" + strconv.Itoa(cur)
	}

	s.render(w, r, "lesson_steps.html", ViewData{Title: lesson.Title, Data: map[string]any{
		"bodyClass":   "lesson-lab",
		"course":      course,
		"lesson":      lesson,
		"step":        step,
		"instruction": content.Render(step.Instruction),
		"checks":      checks,
		"index":       cur + 1,
		"total":       len(steps),
		"progressPct": (cur + 1) * 100 / len(steps),
		"completed":   completed[step.ID],
		"isLast":      isLast,
		"nextURL":     nextURL,
		"prevURL":     prevURL,
		"completeURL": "/steps/" + strconv.FormatInt(step.ID, 10) + "/complete",
		"runURL":      "/steps/" + strconv.FormatInt(step.ID, 10) + "/run",
	}})
}

func (s *Server) handleStepComplete(w http.ResponseWriter, r *http.Request) {
	u := auth.CurrentUser(r.Context())
	if u == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if _, err := s.st.GetStep(r.Context(), id); err != nil {
		s.notFound(w, r)
		return
	}
	if err := s.st.MarkStepComplete(r.Context(), u.ID, id); err != nil {
		if writeWorkAccessError(w, err) {
			return
		}
		http.Error(w, "could not save progress", http.StatusInternalServerError)
		return
	}
	// Finishing the last step finishes the lesson. Without this a learner who
	// worked through every step would still see the lesson as outstanding on
	// their schedule, and would have to confirm it a second time.
	if st, err := s.st.GetStep(r.Context(), id); err == nil {
		_, _ = s.st.MarkLessonCompleteIfStepsDone(r.Context(), u.ID, st.LessonID)
	}
	w.WriteHeader(http.StatusNoContent)
}

// handleLessonComplete marks a reading lesson finished.
func (s *Server) handleLessonComplete(w http.ResponseWriter, r *http.Request) {
	u := auth.CurrentUser(r.Context())
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	lesson, err := s.st.GetLessonByID(r.Context(), id)
	if err != nil {
		s.notFound(w, r)
		return
	}
	if err := s.st.MarkLessonComplete(r.Context(), u.ID, id); err != nil {
		if writeWorkAccessError(w, err) {
			return
		}
		http.Error(w, "could not save progress", http.StatusInternalServerError)
		return
	}
	course, err := s.st.GetCourseByID(r.Context(), lesson.CourseID)
	if err != nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/courses/"+course.Slug+"/"+lesson.Slug, http.StatusSeeOther)
}

// --- Admin ------------------------------------------------------------------

func (s *Server) handleAdminCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := s.st.ListCourses(r.Context(), true)
	if err != nil {
		http.Error(w, "could not load courses", http.StatusInternalServerError)
		return
	}
	s.render(w, r, "admin_courses.html", ViewData{Title: "Courses", Data: map[string]any{"courses": courses}})
}

// handleAdminSoon renders a placeholder for admin sections still being built.
func (s *Server) handleAdminSoon(title string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.render(w, r, "admin_soon.html", ViewData{Title: title, Data: map[string]any{"section": title}})
	}
}

func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	s.render(w, r, "notfound.html", ViewData{Title: "Not found"})
}
