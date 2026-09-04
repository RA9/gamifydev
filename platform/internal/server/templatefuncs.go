package server

import (
	"encoding/json"
	"html/template"
	"strconv"
	"strings"
	"time"
)

// jsonValue embeds a Go value as JSON inside a data- attribute.
//
// Returns template.JS rather than a string so html/template doesn't escape the
// quotes into entities; the value is ours (starter code an admin authored), and
// it is read back with JSON.parse rather than evaluated.
func jsonValue(v any) template.JS {
	b, err := json.Marshal(v)
	if err != nil {
		return template.JS("null")
	}
	return template.JS(b)
}

// pct renders a/b as a whole percentage, and 0 rather than a panic when b is 0.
func pct(a, b int) int {
	if b <= 0 {
		return 0
	}
	return a * 100 / b
}

// initials returns up to two uppercase initials for a display name, used for
// avatar chips in the admin UI.
func initials(name string) string {
	fields := strings.Fields(name)
	if len(fields) == 0 {
		return "?"
	}
	first := []rune(fields[0])
	out := strings.ToUpper(string(first[0]))
	if len(fields) > 1 {
		last := []rune(fields[len(fields)-1])
		out += strings.ToUpper(string(last[0]))
	}
	return out
}

// parseDBTime parses the timestamp formats SQLite hands back for our columns.
func parseDBTime(s string) (time.Time, bool) {
	for _, layout := range []string{"2006-01-02 15:04:05", time.RFC3339, "2006-01-02"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC(), true
		}
	}
	return time.Time{}, false
}

// fmtDate formats a stored timestamp as e.g. "Jun 24, 2026".
func fmtDate(s string) string {
	if t, ok := parseDBTime(s); ok {
		return t.Format("Jan 2, 2006")
	}
	return s
}

// fmtTime renders just the clock part of a stored instant, for standup windows.
// The stored value is UTC; the label around it names the cohort's band, so this
// stays honest without pretending to know the viewer's own clock.
func fmtTime(s string) string {
	if t, ok := parseDBTime(s); ok {
		return t.Format("15:04") + " UTC"
	}
	return s
}

// relTime renders a compact "time ago" string (e.g. "3h ago", "just now").
func relTime(s string) string {
	t, ok := parseDBTime(s)
	if !ok {
		return s
	}
	d := time.Since(t)
	atLeast := func(n int) string {
		if n < 1 {
			n = 1
		}
		return strconv.Itoa(n)
	}
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return atLeast(int(d.Minutes())) + "m ago"
	case d < 24*time.Hour:
		return atLeast(int(d.Hours())) + "h ago"
	case d < 7*24*time.Hour:
		return atLeast(int(d.Hours()/24)) + "d ago"
	default:
		return t.Format("Jan 2")
	}
}
