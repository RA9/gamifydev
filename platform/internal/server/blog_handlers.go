package server

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/content"
	"github.com/RA9/gamifydev/platform/internal/store"
)

// --- Public -----------------------------------------------------------------

func (s *Server) handleBlog(w http.ResponseWriter, r *http.Request) {
	posts, err := s.st.ListPosts(r.Context(), true)
	if err != nil {
		http.Error(w, "could not load the blog", http.StatusInternalServerError)
		return
	}
	s.render(w, r, "blog.html", ViewData{Title: "Blog", Data: map[string]any{"posts": posts}})
}

func (s *Server) handleBlogPost(w http.ResponseWriter, r *http.Request) {
	p, err := s.st.GetPublishedPostBySlug(r.Context(), r.PathValue("slug"))
	if err != nil {
		s.notFound(w, r)
		return
	}
	s.render(w, r, "blog_post.html", ViewData{
		Title: p.Title,
		Data:  map[string]any{"post": p, "body": content.Render(p.Body)},
	})
}

// --- Admin ------------------------------------------------------------------

func (s *Server) handleAdminBlog(w http.ResponseWriter, r *http.Request) {
	posts, err := s.st.ListPosts(r.Context(), false)
	if err != nil {
		http.Error(w, "could not load posts", http.StatusInternalServerError)
		return
	}
	s.render(w, r, "admin_blog.html", ViewData{Title: "Blog", Data: map[string]any{"posts": posts}})
}

func (s *Server) handleAdminBlogNew(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "admin_blog_form.html", ViewData{Title: "New post", Data: map[string]any{"action": "/admin/blog"}})
}

// handleAdminBlogPreview renders submitted markdown to HTML for the live editor
// preview. It uses the same renderer as published posts, so the preview matches
// the final output (including :::tip / :::warning blocks).
func (s *Server) handleAdminBlogPreview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	body := r.FormValue("body")
	if strings.TrimSpace(body) == "" {
		_, _ = w.Write([]byte(`<p class="muted">Nothing to preview yet — start writing on the right.</p>`))
		return
	}
	_, _ = w.Write([]byte(content.Render(body)))
}

func (s *Server) handleAdminBlogCreate(w http.ResponseWriter, r *http.Request) {
	u := auth.CurrentUser(r.Context())
	p := postFromForm(r)
	p.AuthorID = sql.NullInt64{Int64: u.ID, Valid: true}
	if p.Published {
		p.PublishedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if len(p.Title) < 2 {
		s.render(w, r, "admin_blog_form.html", ViewData{Title: "New post", Flash: "A title is required.", Data: map[string]any{"post": &p, "action": "/admin/blog"}})
		return
	}
	if _, err := s.st.CreatePost(r.Context(), p); err != nil {
		s.render(w, r, "admin_blog_form.html", ViewData{Title: "New post", Flash: "Could not create — is the slug unique?", Data: map[string]any{"post": &p, "action": "/admin/blog"}})
		return
	}
	http.Redirect(w, r, "/admin/blog", http.StatusSeeOther)
}

func (s *Server) handleAdminBlogEdit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	p, err := s.st.GetPostByID(r.Context(), id)
	if err != nil {
		s.notFound(w, r)
		return
	}
	s.render(w, r, "admin_blog_form.html", ViewData{Title: "Edit post", Data: map[string]any{"post": p, "action": "/admin/blog/" + r.PathValue("id")}})
}

func (s *Server) handleAdminBlogUpdate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	existing, err := s.st.GetPostByID(r.Context(), id)
	if err != nil {
		s.notFound(w, r)
		return
	}
	p := postFromForm(r)
	// Stamp the publish date the first time it goes live; keep it thereafter.
	if p.Published {
		if existing.PublishedAt != "" {
			p.PublishedAt = existing.PublishedAt
		} else {
			p.PublishedAt = time.Now().UTC().Format(time.RFC3339)
		}
	}
	if err := s.st.UpdatePost(r.Context(), id, p); err != nil {
		http.Error(w, "could not update post", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/blog", http.StatusSeeOther)
}

func (s *Server) handleAdminBlogDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	_ = s.st.DeletePost(r.Context(), id)
	http.Redirect(w, r, "/admin/blog", http.StatusSeeOther)
}

func postFromForm(r *http.Request) store.Post {
	title := strings.TrimSpace(r.FormValue("title"))
	return store.Post{
		Title:     title,
		Slug:      slugOr(r.FormValue("slug"), title),
		Excerpt:   strings.TrimSpace(r.FormValue("excerpt")),
		Body:      r.FormValue("body"),
		CoverURL:  strings.TrimSpace(r.FormValue("cover_url")),
		Published: r.FormValue("published") == "1",
	}
}
