// Package content renders lesson markdown to HTML. On top of standard markdown
// (via goldmark) it understands the GamifyDev fenced blocks the existing course
// content uses: :::tip / :::analogy / :::warning / :::key / :::example /
// :::project callouts, and the :::quiz / :::predict / :::reorder / :::fill /
// :::match practice blocks. Lesson content is authored by the team (trusted),
// so raw HTML is allowed through.
package content

import (
	"bytes"
	"html/template"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	ghtml "github.com/yuin/goldmark/renderer/html"
)

// md renders trusted, team-authored lesson content (raw HTML allowed).
var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithRendererOptions(ghtml.WithUnsafe()),
)

// mdSafe renders untrusted, user-generated content (forum posts): raw HTML is
// escaped, not passed through — no WithUnsafe.
var mdSafe = goldmark.New(goldmark.WithExtensions(extension.GFM))

// RenderSafe converts user-submitted markdown to HTML, escaping any raw HTML to
// prevent XSS. Use this for anything a learner can type.
func RenderSafe(source string) template.HTML {
	var buf bytes.Buffer
	if err := mdSafe.Convert([]byte(source), &buf); err != nil {
		return template.HTML(template.HTMLEscapeString(source))
	}
	return template.HTML(buf.String())
}

var fenceOpen = regexp.MustCompile(`^:::([a-zA-Z]+)\s*$`)

var calloutMeta = map[string]struct{ Icon, Label, Class string }{
	"tip":     {"💡", "Tip", "tip"},
	"analogy": {"🧩", "Think of it like…", "analogy"},
	"warning": {"⚠️", "Watch out", "warning"},
	"key":     {"🔑", "Key idea", "key"},
	"example": {"🧪", "Example", "example"},
	"project": {"🚀", "What you'll build", "project"},
}

var exerciseLabel = map[string]string{
	"quiz":    "Check yourself",
	"predict": "Predict the output",
	"reorder": "Put it in order",
	"fill":    "Fill in the blank",
	"match":   "Match them up",
}

// Render converts lesson markdown (with GamifyDev fences) into safe HTML.
func Render(source string) template.HTML {
	lines := strings.Split(source, "\n")
	var out strings.Builder
	var normal []string

	flushNormal := func() {
		if len(normal) == 0 {
			return
		}
		out.WriteString(toHTML(strings.Join(normal, "\n")))
		normal = normal[:0]
	}

	for i := 0; i < len(lines); i++ {
		m := fenceOpen.FindStringSubmatch(lines[i])
		if m == nil {
			normal = append(normal, lines[i])
			continue
		}
		// Collect the block body until a lone ":::".
		typ := strings.ToLower(m[1])
		var body []string
		i++
		for i < len(lines) && strings.TrimSpace(lines[i]) != ":::" {
			body = append(body, lines[i])
			i++
		}
		flushNormal()
		if c, ok := calloutMeta[typ]; ok {
			out.WriteString(`<aside class="callout callout-` + c.Class + `"><p class="callout-title">` +
				c.Icon + ` ` + c.Label + `</p><div class="callout-body">` +
				toHTML(strings.Join(body, "\n")) + `</div></aside>`)
		} else if label, ok := exerciseLabel[typ]; ok {
			out.WriteString(renderExercise(label, body))
		} else {
			// Unknown fence: render its body plainly.
			out.WriteString(toHTML(strings.Join(body, "\n")))
		}
	}
	flushNormal()
	return template.HTML(out.String())
}

// renderExercise shows a practice block as a readable card: the prompt, the
// options (correct one marked), and the explanation behind a reveal.
func renderExercise(label string, body []string) string {
	var transformed []string
	for _, ln := range body {
		t := strings.TrimSpace(ln)
		switch {
		case strings.HasPrefix(t, "Q:"):
			transformed = append(transformed, "**"+strings.TrimSpace(t[2:])+"**")
		case strings.HasPrefix(t, "E:"):
			transformed = append(transformed, "", "_Why:_ "+strings.TrimSpace(t[2:]))
		case strings.HasPrefix(t, "- "):
			opt := strings.TrimSpace(t[2:])
			if strings.HasSuffix(opt, "*") {
				opt = "✅ " + strings.TrimSpace(strings.TrimSuffix(opt, "*"))
			} else {
				opt = "• " + opt
			}
			transformed = append(transformed, opt+"  ")
		default:
			transformed = append(transformed, ln)
		}
	}
	return `<div class="exercise"><p class="exercise-title">✏️ ` + template.HTMLEscapeString(label) +
		`</p><div class="exercise-body">` + toHTML(strings.Join(transformed, "\n")) + `</div></div>`
}

func toHTML(source string) string {
	if strings.TrimSpace(source) == "" {
		return ""
	}
	var buf bytes.Buffer
	if err := md.Convert([]byte(source), &buf); err != nil {
		return template.HTMLEscapeString(source)
	}
	return buf.String()
}
