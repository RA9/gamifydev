package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/content"
	"github.com/RA9/gamifydev/platform/internal/store"
)

// assignmentRow is an assignment plus the current learner's latest status.
type assignmentRow struct {
	store.Assignment
	Status string // "" if not attempted
}

// --- Learner ----------------------------------------------------------------

func (s *Server) handleAssignments(w http.ResponseWriter, r *http.Request) {
	u := auth.CurrentUser(r.Context())
	list, err := s.st.ListAssignments(r.Context(), false)
	if err != nil {
		http.Error(w, "could not load assignments", http.StatusInternalServerError)
		return
	}
	rows := make([]assignmentRow, 0, len(list))
	for _, a := range list {
		row := assignmentRow{Assignment: a}
		if sub, err := s.st.LatestSubmission(r.Context(), a.ID, u.ID); err == nil {
			row.Status = sub.Status
		}
		rows = append(rows, row)
	}
	s.render(w, r, "assignments.html", ViewData{Title: "Assignments", Data: map[string]any{"rows": rows}})
}

func (s *Server) handleAssignment(w http.ResponseWriter, r *http.Request) {
	u := auth.CurrentUser(r.Context())
	a, err := s.st.GetAssignmentBySlug(r.Context(), r.PathValue("slug"))
	if err != nil {
		s.notFound(w, r)
		return
	}
	data := map[string]any{
		"assignment": a,
		"prompt":     content.Render(a.Prompt),
		"submitted":  r.URL.Query().Get("submitted") == "1",
	}
	if sub, err := s.st.LatestSubmission(r.Context(), a.ID, u.ID); err == nil {
		data["submission"] = sub
		// Automated check results, so a learner sees exactly which requirement
		// failed rather than waiting on a human to tell them.
		if results, err := s.st.SubmissionCheckResults(r.Context(), sub.ID); err == nil && len(results) > 0 {
			data["checks"] = results
		}
	}
	// The files the program will find beside it in the sandbox. A prompt that
	// says "you are given access.log" is unanswerable unless the learner can
	// read access.log, so they are shown here rather than only written at run
	// time.
	if files, err := s.st.ListFiles(r.Context(), a.ID); err == nil && len(files) > 0 {
		data["files"] = files
	}
	// Whether this checkpoint is machine-gradable at all shapes what we promise
	// the learner about how fast they'll hear back.
	auto, _ := s.st.AutoGradable(r.Context(), a.ID)
	data["autoGraded"] = auto && s.exec != nil && s.exec.Enabled()
	s.render(w, r, "assignment.html", ViewData{Title: a.Title, Data: data})
}

func (s *Server) handleSubmitAssignment(w http.ResponseWriter, r *http.Request) {
	u := auth.CurrentUser(r.Context())
	a, err := s.st.GetAssignmentBySlug(r.Context(), r.PathValue("slug"))
	if err != nil {
		s.notFound(w, r)
		return
	}
	code := strings.TrimSpace(r.FormValue("code"))
	note := strings.TrimSpace(r.FormValue("note"))
	if code == "" {
		http.Redirect(w, r, "/assignments/"+a.Slug, http.StatusSeeOther)
		return
	}
	if _, err := s.st.CreateSubmission(r.Context(), a.ID, u.ID, code, note); err != nil {
		http.Error(w, "could not save submission", http.StatusInternalServerError)
		return
	}
	// Refresh the learner's own dashboard (e.g. an open tab) in realtime.
	s.hub.Notify(u.ID, "submissions")
	http.Redirect(w, r, "/assignments/"+a.Slug+"?submitted=1", http.StatusSeeOther)
}

// --- Staff (admin/grader) ---------------------------------------------------

func (s *Server) handleAdminGrading(w http.ResponseWriter, r *http.Request) {
	queue, err := s.st.GradingQueue(r.Context(), false)
	if err != nil {
		http.Error(w, "could not load queue", http.StatusInternalServerError)
		return
	}
	pending := 0
	for _, q := range queue {
		if q.Status != "graded" {
			pending++
		}
	}
	s.render(w, r, "admin_grading.html", ViewData{
		Title: "Grading queue",
		Data:  map[string]any{"queue": queue, "pending": pending},
	})
}

func (s *Server) handleAdminGradeForm(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	sub, err := s.st.GetSubmission(r.Context(), id)
	if err != nil {
		s.notFound(w, r)
		return
	}
	s.render(w, r, "admin_grade.html", ViewData{
		Title: "Review submission",
		Data:  map[string]any{"sub": sub, "prompt": content.Render(sub.Assignment.Prompt)},
	})
}

func (s *Server) handleAdminGrade(w http.ResponseWriter, r *http.Request) {
	u := auth.CurrentUser(r.Context())
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	score, _ := strconv.Atoi(r.FormValue("score"))
	feedback := strings.TrimSpace(r.FormValue("feedback"))
	status := r.FormValue("status")
	if status != "graded" && status != "returned" {
		status = "graded"
	}
	if err := s.st.GradeSubmission(r.Context(), id, u.ID, score, feedback, status); err != nil {
		http.Error(w, "could not save grade", http.StatusInternalServerError)
		return
	}
	// Push a realtime refresh to the learner whose work was just graded.
	if sub, err := s.st.GetSubmission(r.Context(), id); err == nil {
		s.hub.Notify(sub.UserID, "submissions")
	}
	http.Redirect(w, r, "/admin/grading", http.StatusSeeOther)
}
