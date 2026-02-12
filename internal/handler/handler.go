package handler

import (
	"encoding/json"
	"headcontrol/internal/headscale"
	"headcontrol/internal/model"
	"headcontrol/internal/store"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type Handler struct {
	store     *store.Store
	templates *template.Template
}

func New(s *store.Store, templateDir string) (*Handler, error) {
	funcMap := template.FuncMap{
		"join": strings.Join,
		"json": func(v interface{}) template.JS {
			b, _ := json.Marshal(v)
			return template.JS(b)
		},
		"now":          func() time.Time { return time.Now() },
		"fmtTime":      formatTime,
		"fmtTimeShort": formatTimeShort,
		"timeAgo":      timeAgo,
		"nodeUser": func(n model.Node) string {
			if n.User != nil {
				return n.User.Name
			}
			return "—"
		},
		"nodeIPs": func(n model.Node) string {
			if len(n.IPAddresses) == 0 {
				return "—"
			}
			return strings.Join(n.IPAddresses, ", ")
		},
		"firstIP": func(n model.Node) string {
			if len(n.IPAddresses) == 0 {
				return "—"
			}
			return n.IPAddresses[0]
		},
		"safeLen": func(s []string) int {
			if s == nil {
				return 0
			}
			return len(s)
		},
	}

	tmpl := template.New("").Funcs(funcMap)
	for _, p := range []string{"layout", "pages", "partials"} {
		if _, err := tmpl.ParseGlob(filepath.Join(templateDir, p, "*.html")); err != nil {
			log.Printf("parse %s: %v", p, err)
		}
	}

	return &Handler{store: s, templates: tmpl}, nil
}

func (h *Handler) getClient() (*headscale.Client, error) {
	cfg, err := h.store.GetSettings()
	if err != nil || cfg == nil {
		return nil, err
	}
	return headscale.NewClient(cfg.BaseURL, cfg.APIKey), nil
}

func (h *Handler) render(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("template %s: %v", name, err)
		http.Error(w, "Internal Server Error", 500)
	}
}

func (h *Handler) renderFile(w http.ResponseWriter, path string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles(path)
	if err != nil {
		log.Printf("parse %s: %v", path, err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	if err := t.Execute(w, data); err != nil {
		log.Printf("exec %s: %v", path, err)
		http.Error(w, "Internal Server Error", 500)
	}
}

func (h *Handler) isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func (h *Handler) renderPage(w http.ResponseWriter, r *http.Request, page string, data map[string]interface{}) {
	if h.isHTMX(r) {
		h.render(w, page+"-content.html", data)
	} else {
		h.render(w, "layout.html", data)
	}
}

func (h *Handler) renderPageWithError(w http.ResponseWriter, r *http.Request, title, page, msg string) {
	h.renderPage(w, r, page, map[string]interface{}{
		"Title":      title,
		"ActivePage": page,
		"Error":      msg,
	})
}

func (h *Handler) RequireSetup(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !h.store.HasSettings() {
			if h.isHTMX(r) {
				w.Header().Set("HX-Redirect", "/setup")
				w.WriteHeader(200)
			} else {
				http.Redirect(w, r, "/setup", http.StatusFound)
			}
			return
		}
		next(w, r)
	}
}
