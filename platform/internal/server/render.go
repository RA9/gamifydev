package server

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/store"
)

//go:embed web/templates/*.html
var templatesFS embed.FS

//go:embed web/static
var staticFS embed.FS

// StaticHandler serves embedded assets under /static/.
func StaticHandler() http.Handler {
	sub, _ := fs.Sub(staticFS, "web/static")
	return http.StripPrefix("/static/", http.FileServer(http.FS(sub)))
}

// ViewData is the envelope every page receives.
type ViewData struct {
	Title string
	User  *store.User
	Flash string
	Data  map[string]any
}

type renderer struct {
	pages map[string]*template.Template
}

// newRenderer compiles each page template together with the shared layout.
func newRenderer() (*renderer, error) {
	layout := "web/templates/layout.html"
	entries, err := templatesFS.ReadDir("web/templates")
	if err != nil {
		return nil, err
	}
	r := &renderer{pages: map[string]*template.Template{}}
	for _, e := range entries {
		if e.IsDir() || e.Name() == "layout.html" {
			continue
		}
		t, err := template.New("layout.html").Funcs(funcMap).ParseFS(templatesFS, layout, "web/templates/"+e.Name())
		if err != nil {
			return nil, err
		}
		r.pages[e.Name()] = t
	}
	return r, nil
}

var funcMap = template.FuncMap{
	"title": func(s string) string {
		if s == "" {
			return ""
		}
		return string(s[0]-32) + s[1:]
	},
}

// render writes a full page (page template + layout) for the request.
func (s *Server) render(w http.ResponseWriter, req *http.Request, page string, vd ViewData) {
	t, ok := s.rnd.pages[page]
	if !ok {
		http.Error(w, "template not found: "+page, http.StatusInternalServerError)
		return
	}
	vd.User = auth.CurrentUser(req.Context())
	if vd.Data == nil {
		vd.Data = map[string]any{}
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "layout.html", vd); err != nil {
		http.Error(w, "render error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = buf.WriteTo(w)
}

// renderPartial renders a single named template (an htmx fragment) from a
// page's template set — used to swap just part of the page.
func (s *Server) renderPartial(w http.ResponseWriter, _ *http.Request, page, name string, data any) {
	t, ok := s.rnd.pages[page]
	if !ok {
		http.Error(w, "template not found: "+page, http.StatusInternalServerError)
		return
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, name, data); err != nil {
		http.Error(w, "render error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = buf.WriteTo(w)
}
