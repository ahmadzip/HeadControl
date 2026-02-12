package handler

import "net/http"

func (h *Handler) SetupPage(w http.ResponseWriter, r *http.Request) {
	if h.store.HasSettings() {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	h.renderFile(w, "templates/pages/setup.html", nil)
}

func (h *Handler) TestConnection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	baseURL := r.FormValue("base_url")
	apiKey := r.FormValue("api_key")

	if baseURL == "" || apiKey == "" {
		h.render(w, "connection-result.html", map[string]interface{}{
			"Success": false,
			"Message": "Both Base URL and API Key are required.",
		})
		return
	}

	if err := newTempClient(baseURL, apiKey).TestConnection(); err != nil {
		h.render(w, "connection-result.html", map[string]interface{}{
			"Success": false,
			"Message": err.Error(),
		})
		return
	}

	h.render(w, "connection-result.html", map[string]interface{}{
		"Success": true,
		"Message": "Connection successful! Headscale server is reachable.",
	})
}

func (h *Handler) SaveSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	baseURL := r.FormValue("base_url")
	apiKey := r.FormValue("api_key")

	if baseURL == "" || apiKey == "" {
		h.render(w, "connection-result.html", map[string]interface{}{
			"Success": false,
			"Message": "Both Base URL and API Key are required.",
		})
		return
	}

	if err := newTempClient(baseURL, apiKey).TestConnection(); err != nil {
		h.render(w, "connection-result.html", map[string]interface{}{
			"Success": false,
			"Message": "Connection test failed: " + err.Error(),
		})
		return
	}

	if err := h.store.SaveSettings(baseURL, apiKey); err != nil {
		h.render(w, "connection-result.html", map[string]interface{}{
			"Success": false,
			"Message": "Failed to save settings: " + err.Error(),
		})
		return
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(200)
}
