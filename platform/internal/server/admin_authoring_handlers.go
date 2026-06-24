package server

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/RA9/gamifydev/platform/internal/store"
)

func slugOr(slug, title string) string {
	if s := strings.TrimSpace(slug); s != "" {
		return makeSlug(s)
	}
	return makeSlug(title)
}

// --- Courses ----------------------------------------------------------------

func (s *Server) handleAdminCourseNew(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "admin_course_form.html", ViewData{Title: "New course"})
}

func (s *Server) handleAdminCourseCreate(w http.ResponseWriter, r *http.Request) {
	c := courseFromForm(r)
	if len(c.Title) < 2 {
		s.render(w, r, "admin_course_form.html", ViewData{Title: "New course", Flash: "A title is required.", Data: map[string]any{"course": c}})
		return
	}
	id, err := s.st.CreateCourse(r.Context(), c)
	if err != nil {
		s.render(w, r, "admin_course_form.html", ViewData{Title: "New course", Flash: "Could not create — is the slug unique?", Data: map[string]any{"course": c}})
		return
	}
	http.Redirect(w, r, "/admin/courses/"+strconv.FormatInt(id, 10), http.StatusSeeOther)
}

func (s *Server) handleAdminCourseManage(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	c, err := s.st.GetCourseByID(r.Context(), id)
	if err != nil {
		s.notFound(w, r)
		return
	}
	lessons, _ := s.st.ListLessons(r.Context(), c.ID)
	s.render(w, r, "admin_course.html", ViewData{Title: c.Title, Data: map[string]any{"course": c, "lessons": lessons}})
}

func (s *Server) handleAdminCourseUpdate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	c := courseFromForm(r)
	if err := s.st.UpdateCourse(r.Context(), id, c); err != nil {
		http.Error(w, "could not update course", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/courses/"+r.PathValue("id"), http.StatusSeeOther)
}

func courseFromForm(r *http.Request) store.Course {
	title := strings.TrimSpace(r.FormValue("title"))
	return store.Course{
		Title:       title,
		Slug:        slugOr(r.FormValue("slug"), title),
		Emoji:       strings.TrimSpace(r.FormValue("emoji")),
		Tagline:     strings.TrimSpace(r.FormValue("tagline")),
		Description: strings.TrimSpace(r.FormValue("description")),
		Published:   r.FormValue("published") == "1",
	}
}

// --- Lessons ----------------------------------------------------------------

func (s *Server) handleAdminLessonNew(w http.ResponseWriter, r *http.Request) {
	cid, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	c, err := s.st.GetCourseByID(r.Context(), cid)
	if err != nil {
		s.notFound(w, r)
		return
	}
	s.render(w, r, "admin_lesson_form.html", ViewData{Title: "New lesson", Data: map[string]any{"course": c, "action": "/admin/courses/" + r.PathValue("id") + "/lessons"}})
}

func (s *Server) handleAdminLessonCreate(w http.ResponseWriter, r *http.Request) {
	cid, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	c, err := s.st.GetCourseByID(r.Context(), cid)
	if err != nil {
		s.notFound(w, r)
		return
	}
	l := lessonFromForm(r)
	l.CourseID = cid
	l.Sort = s.st.NextLessonSort(r.Context(), cid)
	if len(l.Title) < 2 {
		s.render(w, r, "admin_lesson_form.html", ViewData{Title: "New lesson", Flash: "A title is required.", Data: map[string]any{"course": c, "lesson": &l, "action": "/admin/courses/" + r.PathValue("id") + "/lessons"}})
		return
	}
	if err := s.st.CreateLesson(r.Context(), l); err != nil {
		s.render(w, r, "admin_lesson_form.html", ViewData{Title: "New lesson", Flash: "Could not create — is the slug unique within this course?", Data: map[string]any{"course": c, "lesson": &l, "action": "/admin/courses/" + r.PathValue("id") + "/lessons"}})
		return
	}
	http.Redirect(w, r, "/admin/courses/"+r.PathValue("id"), http.StatusSeeOther)
}

func (s *Server) handleAdminLessonEdit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	l, err := s.st.GetLessonByID(r.Context(), id)
	if err != nil {
		s.notFound(w, r)
		return
	}
	c, _ := s.st.GetCourseByID(r.Context(), l.CourseID)
	s.render(w, r, "admin_lesson_form.html", ViewData{Title: "Edit lesson", Data: map[string]any{"course": c, "lesson": l, "action": "/admin/lessons/" + r.PathValue("id")}})
}

func (s *Server) handleAdminLessonUpdate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	existing, err := s.st.GetLessonByID(r.Context(), id)
	if err != nil {
		s.notFound(w, r)
		return
	}
	l := lessonFromForm(r)
	if err := s.st.UpdateLesson(r.Context(), id, l); err != nil {
		http.Error(w, "could not update lesson", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/courses/"+strconv.FormatInt(existing.CourseID, 10), http.StatusSeeOther)
}

func (s *Server) handleAdminLessonDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	existing, err := s.st.GetLessonByID(r.Context(), id)
	if err != nil {
		s.notFound(w, r)
		return
	}
	_ = s.st.DeleteLesson(r.Context(), id)
	http.Redirect(w, r, "/admin/courses/"+strconv.FormatInt(existing.CourseID, 10), http.StatusSeeOther)
}

func lessonFromForm(r *http.Request) store.Lesson {
	title := strings.TrimSpace(r.FormValue("title"))
	return store.Lesson{
		Title:    title,
		Slug:     slugOr(r.FormValue("slug"), title),
		Summary:  strings.TrimSpace(r.FormValue("summary")),
		Body:     r.FormValue("body"),
		VideoURL: strings.TrimSpace(r.FormValue("video_url")),
		AudioURL: strings.TrimSpace(r.FormValue("audio_url")),
	}
}

// --- Assignments ------------------------------------------------------------

func (s *Server) handleAdminAssignments(w http.ResponseWriter, r *http.Request) {
	list, err := s.st.ListAssignments(r.Context(), true)
	if err != nil {
		http.Error(w, "could not load assignments", http.StatusInternalServerError)
		return
	}
	s.render(w, r, "admin_assignments.html", ViewData{Title: "Assignments", Data: map[string]any{"assignments": list}})
}

func (s *Server) handleAdminAssignmentNew(w http.ResponseWriter, r *http.Request) {
	courses, _ := s.st.ListCourses(r.Context(), true)
	s.render(w, r, "admin_assignment_form.html", ViewData{Title: "New assignment", Data: map[string]any{"courses": courses, "action": "/admin/assignments"}})
}

func (s *Server) handleAdminAssignmentCreate(w http.ResponseWriter, r *http.Request) {
	a := assignmentFromForm(r)
	if len(a.Title) < 2 {
		courses, _ := s.st.ListCourses(r.Context(), true)
		s.render(w, r, "admin_assignment_form.html", ViewData{Title: "New assignment", Flash: "A title is required.", Data: map[string]any{"courses": courses, "assignment": &a, "action": "/admin/assignments"}})
		return
	}
	if err := s.st.CreateAssignment(r.Context(), a); err != nil {
		courses, _ := s.st.ListCourses(r.Context(), true)
		s.render(w, r, "admin_assignment_form.html", ViewData{Title: "New assignment", Flash: "Could not create — is the slug unique?", Data: map[string]any{"courses": courses, "assignment": &a, "action": "/admin/assignments"}})
		return
	}
	http.Redirect(w, r, "/admin/assignments", http.StatusSeeOther)
}

func (s *Server) handleAdminAssignmentEdit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	a, err := s.st.GetAssignmentByID(r.Context(), id)
	if err != nil {
		s.notFound(w, r)
		return
	}
	courses, _ := s.st.ListCourses(r.Context(), true)
	s.render(w, r, "admin_assignment_form.html", ViewData{Title: "Edit assignment", Data: map[string]any{"courses": courses, "assignment": a, "action": "/admin/assignments/" + r.PathValue("id")}})
}

func (s *Server) handleAdminAssignmentUpdate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	a := assignmentFromForm(r)
	if err := s.st.UpdateAssignment(r.Context(), id, a); err != nil {
		http.Error(w, "could not update assignment", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/assignments", http.StatusSeeOther)
}

func assignmentFromForm(r *http.Request) store.Assignment {
	title := strings.TrimSpace(r.FormValue("title"))
	pts, _ := strconv.Atoi(r.FormValue("max_points"))
	if pts <= 0 {
		pts = 100
	}
	lang := strings.TrimSpace(r.FormValue("language"))
	if lang == "" {
		lang = "text"
	}
	var courseID sql.NullInt64
	if cid, err := strconv.ParseInt(r.FormValue("course_id"), 10, 64); err == nil && cid > 0 {
		courseID = sql.NullInt64{Int64: cid, Valid: true}
	}
	return store.Assignment{
		Title:     title,
		Slug:      slugOr(r.FormValue("slug"), title),
		CourseID:  courseID,
		Language:  lang,
		Prompt:    r.FormValue("prompt"),
		Starter:   r.FormValue("starter"),
		MaxPoints: pts,
		Published: r.FormValue("published") == "1",
	}
}
