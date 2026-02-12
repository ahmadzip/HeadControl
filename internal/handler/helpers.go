package handler

import (
	"fmt"
	"headcontrol/internal/headscale"
	"html/template"
	"net/http"
	"time"
)

func (h *Handler) renderPartialError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<div class="error-banner"><i data-lucide="x-circle"></i><span>%s</span></div>`,
		template.HTMLEscapeString(msg))
}

func (h *Handler) renderToast(w http.ResponseWriter, msg, kind string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<div class="toast toast-%s" id="action-toast">%s</div>`,
		kind, template.HTMLEscapeString(msg))
}

func newTempClient(url, key string) *headscale.Client {
	return headscale.NewClient(url, key)
}

func formatTime(s string) string {
	t, ok := parseTime(s)
	if !ok {
		return s
	}
	if t.IsZero() || t.Year() < 2000 {
		return "Never"
	}
	return t.Format("Jan 02, 2006 15:04")
}

func formatTimeShort(s string) string {
	t, ok := parseTime(s)
	if !ok {
		return s
	}
	if t.IsZero() || t.Year() < 2000 {
		return "Never"
	}
	return t.Format("Jan 02, 15:04")
}

func timeAgo(s string) string {
	t, ok := parseTime(s)
	if !ok {
		return s
	}
	if t.IsZero() || t.Year() < 2000 {
		return "Never"
	}

	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "Just now"
	case d < time.Hour:
		if m := int(d.Minutes()); m == 1 {
			return "1 min ago"
		} else {
			return fmt.Sprintf("%d mins ago", m)
		}
	case d < 24*time.Hour:
		if h := int(d.Hours()); h == 1 {
			return "1 hour ago"
		} else {
			return fmt.Sprintf("%d hours ago", h)
		}
	default:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		if days < 30 {
			return fmt.Sprintf("%d days ago", days)
		}
		return t.Format("Jan 02, 2006")
	}
}

func parseTime(s string) (time.Time, bool) {
	if s == "" {
		return time.Time{}, false
	}
	for _, layout := range []string{time.RFC3339, time.RFC3339Nano} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}
