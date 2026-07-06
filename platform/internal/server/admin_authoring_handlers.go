package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/RA9/gamifydev/platform/internal/content"
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
	assignments, _ := s.st.ListAssignmentsByCourse(r.Context(), c.ID, true)
	s.render(w, r, "admin_course.html", ViewData{Title: c.Title, Data: map[string]any{"course": c, "lessons": lessons, "assignments": assignments}})
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
	steps, _ := s.st.ListSteps(r.Context(), l.ID)
	type stepRow struct {
		store.Step
		CheckCount int
	}
	rows := make([]stepRow, 0, len(steps))
	for _, st := range steps {
		var cs []stepCheck
		_ = json.Unmarshal([]byte(st.Checks), &cs)
		rows = append(rows, stepRow{Step: st, CheckCount: len(cs)})
	}
	s.render(w, r, "admin_lesson_form.html", ViewData{Title: "Edit lesson", Data: map[string]any{"course": c, "lesson": l, "steps": rows, "action": "/admin/lessons/" + r.PathValue("id")}})
}

// --- Interactive steps (admin authoring) ------------------------------------

func (s *Server) handleAdminStepNew(w http.ResponseWriter, r *http.Request) {
	lid, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	l, err := s.st.GetLessonByID(r.Context(), lid)
	if err != nil {
		s.notFound(w, r)
		return
	}
	s.render(w, r, "admin_step_form.html", ViewData{Title: "New step", Data: map[string]any{
		"lesson": l, "action": "/admin/lessons/" + r.PathValue("id") + "/steps",
	}})
}

func (s *Server) handleAdminStepCreate(w http.ResponseWriter, r *http.Request) {
	lid, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if _, err := s.st.GetLessonByID(r.Context(), lid); err != nil {
		s.notFound(w, r)
		return
	}
	st := stepFromForm(r)
	st.LessonID = lid
	if _, err := s.st.CreateStep(r.Context(), st); err != nil {
		http.Error(w, "could not create step", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/lessons/"+r.PathValue("id"), http.StatusSeeOther)
}

func (s *Server) handleAdminStepEdit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	st, err := s.st.GetStep(r.Context(), id)
	if err != nil {
		s.notFound(w, r)
		return
	}
	l, _ := s.st.GetLessonByID(r.Context(), st.LessonID)
	var checks []stepCheck
	_ = json.Unmarshal([]byte(st.Checks), &checks)
	s.render(w, r, "admin_step_form.html", ViewData{Title: "Edit step", Data: map[string]any{
		"lesson": l, "step": st, "checks": checks, "action": "/admin/steps/" + r.PathValue("id"),
	}})
}

func (s *Server) handleAdminStepUpdate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	existing, err := s.st.GetStep(r.Context(), id)
	if err != nil {
		s.notFound(w, r)
		return
	}
	if err := s.st.UpdateStep(r.Context(), id, stepFromForm(r)); err != nil {
		http.Error(w, "could not save step", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/lessons/"+strconv.FormatInt(existing.LessonID, 10), http.StatusSeeOther)
}

func (s *Server) handleAdminStepDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	st, err := s.st.GetStep(r.Context(), id)
	if err != nil {
		s.notFound(w, r)
		return
	}
	_ = s.st.DeleteStep(r.Context(), id)
	http.Redirect(w, r, "/admin/lessons/"+strconv.FormatInt(st.LessonID, 10), http.StatusSeeOther)
}

func (s *Server) handleAdminStepMove(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	st, err := s.st.GetStep(r.Context(), id)
	if err != nil {
		s.notFound(w, r)
		return
	}
	dir := -1
	if r.PathValue("dir") == "down" {
		dir = 1
	}
	_ = s.st.MoveStep(r.Context(), id, dir)
	http.Redirect(w, r, "/admin/lessons/"+strconv.FormatInt(st.LessonID, 10), http.StatusSeeOther)
}

// stepFromForm builds a Step from the editor form, zipping the parallel
// check_text[] / check_test[] inputs into the checks JSON.
func stepFromForm(r *http.Request) store.Step {
	_ = r.ParseForm()
	texts := r.Form["check_text"]
	tests := r.Form["check_test"]
	var checks []stepCheck
	for i := range texts {
		text := strings.TrimSpace(texts[i])
		test := ""
		if i < len(tests) {
			test = strings.TrimSpace(tests[i])
		}
		if text == "" && test == "" {
			continue
		}
		checks = append(checks, stepCheck{Text: text, Test: test})
	}
	if checks == nil {
		checks = []stepCheck{}
	}
	raw, _ := json.Marshal(checks)
	lang := "html"
	switch r.FormValue("lang") {
	case "js":
		lang = "js"
	case "python":
		lang = "python"
	}
	return store.Step{
		Instruction: r.FormValue("instruction"),
		Starter:     r.FormValue("starter"),
		Checks:      string(raw),
		Lang:        lang,
		Scaffold:    r.FormValue("scaffold"),
	}
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

func (s *Server) handleAdminLessonPreview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	body := r.FormValue("body")
	if strings.TrimSpace(body) == "" {
		_, _ = w.Write([]byte(`<p class="muted">Nothing to preview yet — start writing.</p>`))
		return
	}
	_, _ = w.Write([]byte(content.Render(body)))
}

func lessonFromForm(r *http.Request) store.Lesson {
	title := strings.TrimSpace(r.FormValue("title"))
	kind := strings.TrimSpace(r.FormValue("kind"))
	if kind == "" {
		kind = "theory"
	}
	return store.Lesson{
		Title:    title,
		Slug:     slugOr(r.FormValue("slug"), title),
		Summary:  strings.TrimSpace(r.FormValue("summary")),
		Body:     r.FormValue("body"),
		VideoURL: strings.TrimSpace(r.FormValue("video_url")),
		AudioURL: strings.TrimSpace(r.FormValue("audio_url")),
		Section:  strings.TrimSpace(r.FormValue("section")),
		Kind:     kind,
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
	// Return to the owning course's manager when the assignment belongs to one.
	if a.CourseID.Valid {
		http.Redirect(w, r, "/admin/courses/"+strconv.FormatInt(a.CourseID.Int64, 10), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/admin/assignments", http.StatusSeeOther)
}

// --- Assignments scoped to a course (add directly from the course manager) ---

func (s *Server) handleAdminCourseAssignmentNew(w http.ResponseWriter, r *http.Request) {
	cid, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	c, err := s.st.GetCourseByID(r.Context(), cid)
	if err != nil {
		s.notFound(w, r)
		return
	}
	courses, _ := s.st.ListCourses(r.Context(), true)
	back := "/admin/courses/" + r.PathValue("id")
	s.render(w, r, "admin_assignment_form.html", ViewData{Title: "New assignment", Data: map[string]any{
		"courses": courses, "presetCourseID": c.ID, "action": back + "/assignments", "back": back,
	}})
}

func (s *Server) handleAdminCourseAssignmentCreate(w http.ResponseWriter, r *http.Request) {
	cid, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	c, err := s.st.GetCourseByID(r.Context(), cid)
	if err != nil {
		s.notFound(w, r)
		return
	}
	back := "/admin/courses/" + r.PathValue("id")
	a := assignmentFromForm(r)
	a.CourseID = sql.NullInt64{Int64: cid, Valid: true} // force this course
	existing, _ := s.st.ListAssignmentsByCourse(r.Context(), cid, true)
	a.Sort = len(existing)
	if len(a.Title) < 2 {
		courses, _ := s.st.ListCourses(r.Context(), true)
		s.render(w, r, "admin_assignment_form.html", ViewData{Title: "New assignment", Flash: "A title is required.",
			Data: map[string]any{"courses": courses, "assignment": &a, "presetCourseID": c.ID, "action": back + "/assignments", "back": back}})
		return
	}
	if err := s.st.CreateAssignment(r.Context(), a); err != nil {
		courses, _ := s.st.ListCourses(r.Context(), true)
		s.render(w, r, "admin_assignment_form.html", ViewData{Title: "New assignment", Flash: "Could not create — is the slug unique?",
			Data: map[string]any{"courses": courses, "assignment": &a, "presetCourseID": c.ID, "action": back + "/assignments", "back": back}})
		return
	}
	http.Redirect(w, r, back, http.StatusSeeOther)
}

func (s *Server) handleAdminAssignmentDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	a, err := s.st.GetAssignmentByID(r.Context(), id)
	if err != nil {
		s.notFound(w, r)
		return
	}
	_ = s.st.DeleteAssignment(r.Context(), id)
	if a.CourseID.Valid {
		http.Redirect(w, r, "/admin/courses/"+strconv.FormatInt(a.CourseID.Int64, 10), http.StatusSeeOther)
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
	passPts, _ := strconv.Atoi(r.FormValue("pass_points"))
	if passPts < 0 {
		passPts = 0
	}
	return store.Assignment{
		Title:      title,
		Slug:       slugOr(r.FormValue("slug"), title),
		CourseID:   courseID,
		Language:   lang,
		Prompt:     r.FormValue("prompt"),
		Starter:    r.FormValue("starter"),
		MaxPoints:  pts,
		Published:  r.FormValue("published") == "1",
		Required:   r.FormValue("required") == "1",
		PassPoints: passPts,
	}
}
