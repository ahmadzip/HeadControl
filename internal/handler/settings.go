package handler

import (
	"headcontrol/internal/model"
	"net/http"
)

func (h *Handler) SettingsPage(w http.ResponseWriter, r *http.Request) {
	settings, _ := h.store.GetSettings()

	var masked *model.Settings
	if settings != nil {
		masked = &model.Settings{
			ID:        settings.ID,
			BaseURL:   settings.BaseURL,
			CreatedAt: settings.CreatedAt,
			UpdatedAt: settings.UpdatedAt,
		}
	}

	h.renderPage(w, r, "settings", map[string]interface{}{
		"Title":      "Settings",
		"ActivePage": "settings",
		"Settings":   masked,
	})
}

func (h *Handler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	baseURL := r.FormValue("base_url")
	apiKey := r.FormValue("api_key")

	if baseURL == "" {
		h.render(w, "settings-result.html", map[string]interface{}{
			"Success": false,
			"Message": "Base URL is required.",
		})
		return
	}

	if apiKey == "" {
		if existing, _ := h.store.GetSettings(); existing != nil {
			apiKey = existing.APIKey
		}
	}

	if apiKey == "" {
		h.render(w, "settings-result.html", map[string]interface{}{
			"Success": false,
			"Message": "API Key is required.",
		})
		return
	}

	if err := h.store.SaveSettings(baseURL, apiKey); err != nil {
		h.render(w, "settings-result.html", map[string]interface{}{
			"Success": false,
			"Message": "Failed to save: " + err.Error(),
		})
		return
	}

	h.render(w, "settings-result.html", map[string]interface{}{
		"Success": true,
		"Message": "Settings saved successfully!",
	})
}
