package server

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/RA9/gamifydev/platform/internal/store"
)

// areaChartSVG renders a signup time-series as a responsive area+line chart.
// The SVG uses a fixed viewBox and scales to its container via CSS; colors come
// from the stylesheet (currentColor / class hooks) so it themes consistently.
func areaChartSVG(series []store.DayCount) template.HTML {
	const w, h = 720.0, 240.0
	const padL, padR, padT, padB = 14.0, 14.0, 16.0, 30.0
	iw, ih := w-padL-padR, h-padT-padB

	n := len(series)
	if n == 0 {
		return ""
	}
	max := 1
	for _, d := range series {
		if d.Count > max {
			max = d.Count
		}
	}
	// Headroom so the peak doesn't touch the top edge.
	scaleMax := float64(max) * 1.15

	xAt := func(i int) float64 {
		if n == 1 {
			return padL + iw/2
		}
		return padL + float64(i)*iw/float64(n-1)
	}
	yAt := func(v int) float64 {
		return padT + ih - (float64(v)/scaleMax)*ih
	}

	var line strings.Builder // the stroke path
	var area strings.Builder // the filled area path
	var dots strings.Builder // point markers
	baseline := padT + ih
	area.WriteString(fmt.Sprintf("M%.1f,%.1f", xAt(0), baseline))
	for i, d := range series {
		x, y := xAt(i), yAt(d.Count)
		if i == 0 {
			line.WriteString(fmt.Sprintf("M%.1f,%.1f", x, y))
		} else {
			line.WriteString(fmt.Sprintf(" L%.1f,%.1f", x, y))
		}
		area.WriteString(fmt.Sprintf(" L%.1f,%.1f", x, y))
	}
	area.WriteString(fmt.Sprintf(" L%.1f,%.1f Z", xAt(n-1), baseline))

	// Emphasised marker on the most recent point.
	lastX, lastY := xAt(n-1), yAt(series[n-1].Count)
	dots.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="5.5" class="chart-dot-halo"/>`, lastX, lastY))
	dots.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="3.5" class="chart-dot"/>`, lastX, lastY))

	// Faint horizontal gridlines (0, mid, top) with value labels.
	var grid strings.Builder
	for _, frac := range []float64{0, 0.5, 1} {
		y := padT + ih - frac*ih
		grid.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" class="chart-grid"/>`, padL, y, w-padR, y))
	}

	// Sparse x-axis labels (first, middle, last) to avoid clutter.
	var labels strings.Builder
	idxs := []int{0, n / 2, n - 1}
	anchors := []string{"start", "middle", "end"}
	for k, i := range idxs {
		if i < 0 || i >= n {
			continue
		}
		labels.WriteString(fmt.Sprintf(
			`<text x="%.1f" y="%.1f" text-anchor="%s" class="chart-xlabel">%s</text>`,
			xAt(i), h-8, anchors[k], template.HTMLEscapeString(series[i].Label)))
	}

	svg := fmt.Sprintf(`<svg viewBox="0 0 %.0f %.0f" class="chart-area" role="img" aria-label="New members per day">
<defs><linearGradient id="areaFill" x1="0" y1="0" x2="0" y2="1">
<stop offset="0%%" class="chart-grad-top"/><stop offset="100%%" class="chart-grad-bottom"/>
</linearGradient></defs>
%s<path d="%s" fill="url(#areaFill)"/><path d="%s" class="chart-line" fill="none"/>%s%s</svg>`,
		w, h, grid.String(), area.String(), line.String(), dots.String(), labels.String())
	return template.HTML(svg)
}

// DonutSegment is one slice of the role donut.
type DonutSegment struct {
	Label string
	Count int
	Class string // CSS class providing the stroke color
	Pct   int
}

// donutSVG renders a distribution donut with a centered total and caption.
// Segments with zero count are skipped.
func donutSVG(total int, caption string, segs []DonutSegment) template.HTML {
	const size = 160.0
	const r = 62.0
	const sw = 22.0
	cx, cy := size/2, size/2
	circ := 2 * 3.141592653589793 * r

	var arcs strings.Builder
	// Track background ring.
	arcs.WriteString(fmt.Sprintf(
		`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="none" stroke-width="%.1f" class="donut-track"/>`,
		cx, cy, r, sw))

	offset := 0.0
	if total > 0 {
		for _, s := range segs {
			if s.Count <= 0 {
				continue
			}
			frac := float64(s.Count) / float64(total)
			dash := frac * circ
			arcs.WriteString(fmt.Sprintf(
				`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="none" stroke-width="%.1f" class="%s" stroke-linecap="round" stroke-dasharray="%.2f %.2f" stroke-dashoffset="%.2f" transform="rotate(-90 %.1f %.1f)"/>`,
				cx, cy, r, sw, s.Class, dash, circ-dash, -offset, cx, cy))
			offset += dash
		}
	}

	svg := fmt.Sprintf(`<svg viewBox="0 0 %.0f %.0f" class="chart-donut" role="img" aria-label="%s">
%s<text x="%.1f" y="%.1f" text-anchor="middle" class="donut-total">%d</text>
<text x="%.1f" y="%.1f" text-anchor="middle" class="donut-sub">%s</text></svg>`,
		size, size, template.HTMLEscapeString(caption), arcs.String(), cx, cy-2, total, cx, cy+18,
		template.HTMLEscapeString(caption))
	return template.HTML(svg)
}
