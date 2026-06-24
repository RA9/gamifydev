package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/content"
)

var forumCategories = []string{"General", "Help", "Show & Tell", "Feedback"}

// --- Public / learner -------------------------------------------------------

func (s *Server) handleForum(w http.ResponseWriter, r *http.Request) {
	threads, err := s.st.ListThreads(r.Context())
	if err != nil {
		http.Error(w, "could not load the forum", http.StatusInternalServerError)
		return
	}
	s.render(w, r, "forum.html", ViewData{Title: "Forum", Data: map[string]any{"threads": threads}})
}

func (s *Server) handleForumNewForm(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "forum_new.html", ViewData{Title: "New discussion", Data: map[string]any{"categories": forumCategories}})
}

func (s *Server) handleForumCreate(w http.ResponseWriter, r *http.Request) {
	u := auth.CurrentUser(r.Context())
	title := strings.TrimSpace(r.FormValue("title"))
	body := strings.TrimSpace(r.FormValue("body"))
	category := r.FormValue("category")
	if !validCategory(category) {
		category = "General"
	}
	if len(title) < 4 || body == "" {
		s.render(w, r, "forum_new.html", ViewData{
			Title: "New discussion", Flash: "Give your post a title (4+ chars) and a body.",
			Data: map[string]any{"categories": forumCategories, "title": title, "body": body},
		})
		return
	}
	id, err := s.st.CreateThread(r.Context(), u.ID, title, body, category)
	if err != nil {
		http.Error(w, "could not create the discussion", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/forum/"+strconv.FormatInt(id, 10), http.StatusSeeOther)
}

func (s *Server) handleForumThread(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	t, err := s.st.GetThread(r.Context(), id)
	if err != nil {
		s.notFound(w, r)
		return
	}
	replies, _ := s.st.ListReplies(r.Context(), id)
	type replyView struct {
		AuthorName, AuthorRole, CreatedAt string
		Body                              any
	}
	rvs := make([]replyView, 0, len(replies))
	for _, rp := range replies {
		rvs = append(rvs, replyView{rp.AuthorName, rp.AuthorRole, rp.CreatedAt, content.RenderSafe(rp.Body)})
	}
	s.render(w, r, "forum_thread.html", ViewData{
		Title: t.Title,
		Data: map[string]any{
			"thread":  t,
			"body":    content.RenderSafe(t.Body),
			"replies": rvs,
		},
	})
}

func (s *Server) handleForumReply(w http.ResponseWriter, r *http.Request) {
	u := auth.CurrentUser(r.Context())
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	t, err := s.st.GetThread(r.Context(), id)
	if err != nil {
		s.notFound(w, r)
		return
	}
	if t.Locked {
		http.Redirect(w, r, "/forum/"+r.PathValue("id"), http.StatusSeeOther)
		return
	}
	body := strings.TrimSpace(r.FormValue("body"))
	if body != "" {
		_ = s.st.CreateReply(r.Context(), id, u.ID, body)
	}
	http.Redirect(w, r, "/forum/"+r.PathValue("id")+"#replies", http.StatusSeeOther)
}

// --- Admin moderation -------------------------------------------------------

func (s *Server) handleAdminForum(w http.ResponseWriter, r *http.Request) {
	threads, err := s.st.ListThreads(r.Context())
	if err != nil {
		http.Error(w, "could not load threads", http.StatusInternalServerError)
		return
	}
	s.render(w, r, "admin_forum.html", ViewData{Title: "Forum", Data: map[string]any{"threads": threads}})
}

func (s *Server) handleForumModerate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	switch r.PathValue("action") {
	case "pin":
		_ = s.st.SetThreadPinned(r.Context(), id, r.FormValue("on") == "1")
	case "lock":
		_ = s.st.SetThreadLocked(r.Context(), id, r.FormValue("on") == "1")
	case "delete":
		_ = s.st.DeleteThread(r.Context(), id)
	}
	http.Redirect(w, r, "/admin/forum", http.StatusSeeOther)
}

func validCategory(c string) bool {
	for _, x := range forumCategories {
		if x == c {
			return true
		}
	}
	return false
}
