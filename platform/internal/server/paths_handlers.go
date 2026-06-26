package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/RA9/gamifydev/platform/internal/content"
	"github.com/RA9/gamifydev/platform/internal/store"
)

// --- Public -----------------------------------------------------------------

func (s *Server) handlePaths(w http.ResponseWriter, r *http.Request) {
	paths, err := s.st.ListPaths(r.Context(), false)
	if err != nil {
		http.Error(w, "could not load career paths", http.StatusInternalServerError)
		return
	}
	s.render(w, r, "paths.html", ViewData{Title: "Career paths", Data: map[string]any{"paths": paths}})
}

func (s *Server) handlePath(w http.ResponseWriter, r *http.Request) {
	p, err := s.st.GetPathBySlug(r.Context(), r.PathValue("slug"))
	if err != nil || !p.Published {
		s.notFound(w, r)
		return
	}
	courses, _ := s.st.PathCourses(r.Context(), p.ID)
	s.render(w, r, "path.html", ViewData{
		Title: p.Title,
		Data: map[string]any{
			"path":    p,
			"courses": courses,
			"about":   content.RenderSafe(p.Description),
		},
	})
}

// --- Admin ------------------------------------------------------------------

func (s *Server) handleAdminPaths(w http.ResponseWriter, r *http.Request) {
	paths, err := s.st.ListPaths(r.Context(), true)
	if err != nil {
		http.Error(w, "could not load paths", http.StatusInternalServerError)
		return
	}
	s.render(w, r, "admin_paths.html", ViewData{Title: "Career paths", Data: map[string]any{"paths": paths}})
}

func (s *Server) handleAdminPathNew(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "admin_path_form.html", ViewData{Title: "New career path"})
}

func (s *Server) handleAdminPathCreate(w http.ResponseWriter, r *http.Request) {
	p := pathFromForm(r)
	if p.Title == "" {
		s.render(w, r, "admin_path_form.html", ViewData{Title: "New career path",
			Flash: "A title is required.", Data: map[string]any{"path": &p}})
		return
	}
	id, err := s.st.CreatePath(r.Context(), p)
	if err != nil {
		s.render(w, r, "admin_path_form.html", ViewData{Title: "New career path",
			Flash: "Could not create — is the slug unique?", Data: map[string]any{"path": &p}})
		return
	}
	// Send to the manage page to assign courses.
	http.Redirect(w, r, "/admin/paths/"+strconv.FormatInt(id, 10), http.StatusSeeOther)
}

func (s *Server) handleAdminPathEdit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	p, err := s.st.GetPathByID(r.Context(), id)
	if err != nil {
		s.notFound(w, r)
		return
	}
	courses, _ := s.st.ListCourses(r.Context(), true)
	assigned, _ := s.st.PathCourseIDs(r.Context(), id)

	// Build display rows: courses already in the path (in order) first, then the rest.
	type row struct {
		Course   store.Course
		Selected bool
		Order    int
	}
	rows := make([]row, 0, len(courses))
	for _, c := range courses {
		sort, ok := assigned[c.ID]
		rows = append(rows, row{Course: c, Selected: ok, Order: sort})
	}
	s.render(w, r, "admin_path_form.html", ViewData{Title: "Edit career path",
		Data: map[string]any{"path": p, "courseRows": rows}})
}

func (s *Server) handleAdminPathUpdate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	p := pathFromForm(r)
	if err := s.st.UpdatePath(r.Context(), id, p); err != nil {
		http.Error(w, "could not save path", http.StatusInternalServerError)
		return
	}
	// Course assignment: checked courses ordered by their "order_<id>" value.
	_ = r.ParseForm()
	type oc struct {
		id    int64
		order int
	}
	var picked []oc
	for _, v := range r.Form["course"] {
		cid, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}
		ord, _ := strconv.Atoi(r.FormValue("order_" + v))
		picked = append(picked, oc{id: cid, order: ord})
	}
	// Stable sort by order.
	for i := 1; i < len(picked); i++ {
		for j := i; j > 0 && picked[j-1].order > picked[j].order; j-- {
			picked[j-1], picked[j] = picked[j], picked[j-1]
		}
	}
	ids := make([]int64, len(picked))
	for i, p := range picked {
		ids[i] = p.id
	}
	if err := s.st.SetPathCourses(r.Context(), id, ids); err != nil {
		http.Error(w, "could not save courses", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/paths", http.StatusSeeOther)
}

func (s *Server) handleAdminPathDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err := s.st.DeletePath(r.Context(), id); err != nil {
		http.Error(w, "could not delete path", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/paths", http.StatusSeeOther)
}

func pathFromForm(r *http.Request) store.Path {
	title := strings.TrimSpace(r.FormValue("title"))
	level := strings.TrimSpace(r.FormValue("level"))
	if level == "" {
		level = "Beginner"
	}
	return store.Path{
		Slug:        slugOr(r.FormValue("slug"), title),
		Title:       title,
		Tagline:     strings.TrimSpace(r.FormValue("tagline")),
		Description: strings.TrimSpace(r.FormValue("description")),
		Emoji:       strings.TrimSpace(r.FormValue("emoji")),
		Level:       level,
		Published:   r.FormValue("published") == "1",
	}
}
