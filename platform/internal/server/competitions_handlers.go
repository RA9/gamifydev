package server

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/content"
	"github.com/RA9/gamifydev/platform/internal/store"
)

var slugRe = regexp.MustCompile(`[^a-z0-9]+`)

func makeSlug(s string) string {
	return strings.Trim(slugRe.ReplaceAllString(strings.ToLower(s), "-"), "-")
}

// --- Public / learner -------------------------------------------------------

func (s *Server) handleCompetitions(w http.ResponseWriter, r *http.Request) {
	comps, err := s.st.ListCompetitions(r.Context(), false)
	if err != nil {
		http.Error(w, "could not load competitions", http.StatusInternalServerError)
		return
	}
	s.render(w, r, "competitions.html", ViewData{Title: "Competitions", Data: map[string]any{"comps": comps}})
}

func (s *Server) handleCompetition(w http.ResponseWriter, r *http.Request) {
	c, err := s.st.GetCompetitionBySlug(r.Context(), r.PathValue("slug"))
	if err != nil {
		s.notFound(w, r)
		return
	}
	board, _ := s.st.Leaderboard(r.Context(), c.ID)
	data := map[string]any{
		"comp":   c,
		"prompt": content.Render(c.Prompt),
		"board":  board,
		"status": c.Status(),
	}
	if u := auth.CurrentUser(r.Context()); u != nil {
		if e, err := s.st.GetEntry(r.Context(), c.ID, u.ID); err == nil {
			data["entry"] = e
		}
	}
	s.render(w, r, "competition.html", ViewData{Title: c.Title, Data: data})
}

func (s *Server) handleCompetitionSubmit(w http.ResponseWriter, r *http.Request) {
	u := auth.CurrentUser(r.Context())
	c, err := s.st.GetCompetitionBySlug(r.Context(), r.PathValue("slug"))
	if err != nil {
		s.notFound(w, r)
		return
	}
	if !c.IsLive() {
		http.Redirect(w, r, "/competitions/"+c.Slug, http.StatusSeeOther)
		return
	}
	code := strings.TrimSpace(r.FormValue("code"))
	if code != "" {
		_ = s.st.UpsertEntry(r.Context(), c.ID, u.ID, code)
	}
	http.Redirect(w, r, "/competitions/"+c.Slug+"?entered=1", http.StatusSeeOther)
}

// --- Admin ------------------------------------------------------------------

func (s *Server) handleAdminCompetitions(w http.ResponseWriter, r *http.Request) {
	comps, err := s.st.ListCompetitions(r.Context(), true)
	if err != nil {
		http.Error(w, "could not load competitions", http.StatusInternalServerError)
		return
	}
	s.render(w, r, "admin_competitions.html", ViewData{Title: "Competitions", Data: map[string]any{"comps": comps}})
}

func (s *Server) handleAdminCompetitionNew(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "admin_competition_new.html", ViewData{Title: "New competition"})
}

func (s *Server) handleAdminCompetitionCreate(w http.ResponseWriter, r *http.Request) {
	title := strings.TrimSpace(r.FormValue("title"))
	prompt := strings.TrimSpace(r.FormValue("prompt"))
	lang := strings.TrimSpace(r.FormValue("language"))
	points, _ := strconv.Atoi(r.FormValue("points"))
	starts, errS := time.ParseInLocation("2006-01-02T15:04", r.FormValue("starts_at"), time.Local)
	ends, errE := time.ParseInLocation("2006-01-02T15:04", r.FormValue("ends_at"), time.Local)

	if len(title) < 4 || errS != nil || errE != nil || !ends.After(starts) {
		s.render(w, r, "admin_competition_new.html", ViewData{
			Title: "New competition",
			Flash: "Give it a title and a valid start/end window (end after start).",
			Data:  map[string]any{"title": title, "prompt": prompt},
		})
		return
	}
	if points <= 0 {
		points = 100
	}
	if lang == "" {
		lang = "text"
	}
	err := s.st.CreateCompetition(r.Context(), store.Competition{
		Slug: makeSlug(title), Title: title, Prompt: prompt, Language: lang,
		Points: points, StartsAt: starts, EndsAt: ends, Published: true,
	})
	if err != nil {
		s.render(w, r, "admin_competition_new.html", ViewData{
			Title: "New competition", Flash: "Could not create (is the title unique?).",
			Data: map[string]any{"title": title, "prompt": prompt},
		})
		return
	}
	http.Redirect(w, r, "/admin/competitions", http.StatusSeeOther)
}

func (s *Server) handleAdminCompetitionManage(w http.ResponseWriter, r *http.Request) {
	c, err := s.st.GetCompetitionBySlug(r.Context(), r.PathValue("slug"))
	if err != nil {
		s.notFound(w, r)
		return
	}
	entries, _ := s.st.ListEntriesWithCode(r.Context(), c.ID)
	s.render(w, r, "admin_competition.html", ViewData{
		Title: c.Title, Data: map[string]any{"comp": c, "entries": entries, "status": c.Status()},
	})
}

func (s *Server) handleAdminScoreEntry(w http.ResponseWriter, r *http.Request) {
	entryID, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	score, _ := strconv.Atoi(r.FormValue("score"))
	_ = s.st.SetEntryScore(r.Context(), entryID, score)
	http.Redirect(w, r, "/admin/competitions/"+r.PathValue("slug"), http.StatusSeeOther)
}
