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

// The three shells. Pages are compiled against all of them and the right one
// is chosen per request:
//   - layout.html        : public / marketing (guests)
//   - app_layout.html    : the learner portal (authenticated, app-style nav)
//   - admin_layout.html  : the admin sidebar
const (
	layoutPublic  = "layout.html"
	layoutApp     = "app_layout.html"
	layoutAdmin   = "admin_layout.html"
	layoutLanding = "landing_layout.html"
)

var allLayouts = []string{layoutPublic, layoutApp, layoutAdmin, layoutLanding}

// alwaysPublic pages keep the marketing shell even when signed in.
var alwaysPublic = map[string]bool{
	"login.html": true, "register.html": true, "notfound.html": true,
}

type renderer struct {
	pages map[string]*template.Template
}

// newRenderer compiles each page together with all layout shells, so the layout
// can be selected at render time (guest vs learner vs admin).
func newRenderer() (*renderer, error) {
	entries, err := templatesFS.ReadDir("web/templates")
	if err != nil {
		return nil, err
	}
	layoutPaths := make([]string, len(allLayouts))
	for i, l := range allLayouts {
		layoutPaths[i] = "web/templates/" + l
	}
	r := &renderer{pages: map[string]*template.Template{}}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || name == layoutPublic || name == layoutApp || name == layoutAdmin {
			continue
		}
		files := append(append([]string{}, layoutPaths...), "web/templates/"+name)
		t, err := template.New(name).Funcs(funcMap).ParseFS(templatesFS, files...)
		if err != nil {
			return nil, err
		}
		r.pages[name] = t
	}
	return r, nil
}

// chooseLayout decides which shell wraps a page for this request.
func chooseLayout(page string, signedIn bool) string {
	switch {
	case page == "home.html":
		return layoutLanding
	case strings.HasPrefix(page, "admin_"):
		return layoutAdmin
	case alwaysPublic[page]:
		return layoutPublic
	case signedIn:
		return layoutApp
	default:
		return layoutPublic
	}
}

var funcMap = template.FuncMap{
	"add":      func(a, b int) int { return a + b },
	"sub":      func(a, b int) int { return a - b },
	"icon":     icon,
	"initials": initials,
	"fmtdate":  fmtDate,
	"reltime":  relTime,
	"titlecase": func(s string) string {
		if s == "" {
			return s
		}
		return strings.ToUpper(s[:1]) + s[1:]
	},
	"hasprefix": strings.HasPrefix,
}

// render writes a full page (page template + layout) for the request.
func (s *Server) render(w http.ResponseWriter, req *http.Request, page string, vd ViewData) {
	t, ok := s.rnd.pages[page]
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
	if err := t.ExecuteTemplate(&buf, chooseLayout(page, vd.User != nil), vd); err != nil {
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
