package server

import (
	"net/http"

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
	assignments, _ := s.st.ListAssignmentsByCourse(r.Context(), course.ID)
	locked := auth.CurrentUser(r.Context()) == nil
	s.render(w, r, "course.html", ViewData{
		Title: course.Title,
		Data: map[string]any{
			"course":      course,
			"lessons":     lessons,
			"assignments": assignments,
			"locked":      locked,
		},
	})
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
	// Guests can see the lesson exists, but the body is gated behind sign-in.
	locked := auth.CurrentUser(r.Context()) == nil
	data := map[string]any{
		"course": course,
		"lesson": lesson,
		"prev":   prev,
		"next":   next,
		"locked": locked,
	}
	if !locked {
		data["body"] = content.Render(lesson.Body)
	}
	s.render(w, r, "lesson.html", ViewData{Title: lesson.Title, Data: data})
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
