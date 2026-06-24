package server

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

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
	Path  string // current request path (for nav active state)
	Data  map[string]any
}

type pageTmpl struct {
	t      *template.Template
	layout string // root template name to execute
}

type renderer struct {
	pages map[string]pageTmpl
}

// newRenderer compiles each page with its layout: admin_* pages use the sidebar
// layout (admin_layout.html); everything else uses the top-nav layout.
func newRenderer() (*renderer, error) {
	entries, err := templatesFS.ReadDir("web/templates")
	if err != nil {
		return nil, err
	}
	r := &renderer{pages: map[string]pageTmpl{}}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || name == "layout.html" || name == "admin_layout.html" {
			continue
		}
		layout := "layout.html"
		if strings.HasPrefix(name, "admin_") {
			layout = "admin_layout.html"
		}
		t, err := template.New(layout).Funcs(funcMap).
			ParseFS(templatesFS, "web/templates/"+layout, "web/templates/"+name)
		if err != nil {
			return nil, err
		}
		r.pages[name] = pageTmpl{t: t, layout: layout}
	}
	return r, nil
}

var funcMap = template.FuncMap{
	"add": func(a, b int) int { return a + b },
}

// render writes a full page (page template + layout) for the request.
func (s *Server) render(w http.ResponseWriter, req *http.Request, page string, vd ViewData) {
	p, ok := s.rnd.pages[page]
	if !ok {
		http.Error(w, "template not found: "+page, http.StatusInternalServerError)
		return
	}
	vd.User = auth.CurrentUser(req.Context())
	vd.Path = req.URL.Path
	if vd.Data == nil {
		vd.Data = map[string]any{}
	}
	var buf bytes.Buffer
	if err := p.t.ExecuteTemplate(&buf, p.layout, vd); err != nil {
		http.Error(w, "render error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = buf.WriteTo(w)
}

// renderPartial renders a single named template (an htmx fragment) from a
// page's template set — used to swap just part of the page.
func (s *Server) renderPartial(w http.ResponseWriter, _ *http.Request, page, name string, data any) {
	p, ok := s.rnd.pages[page]
	if !ok {
		http.Error(w, "template not found: "+page, http.StatusInternalServerError)
		return
	}
	var buf bytes.Buffer
	if err := p.t.ExecuteTemplate(&buf, name, data); err != nil {
		http.Error(w, "render error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = buf.WriteTo(w)
}
