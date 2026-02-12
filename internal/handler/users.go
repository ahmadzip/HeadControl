package handler

import (
	"net/http"
)

func (h *Handler) UsersPage(w http.ResponseWriter, r *http.Request) {
	client, err := h.getClient()
	if err != nil || client == nil {
		h.renderPageWithError(w, r, "Users", "users", "Failed to load settings.")
		return
	}

	users, apiErr := client.ListUsers()
	if apiErr != nil {
		h.renderPageWithError(w, r, "Users", "users", apiErr.Error())
		return
	}

	h.renderPage(w, r, "users", map[string]interface{}{
		"Title":      "Users",
		"ActivePage": "users",
		"Users":      users,
	})
}

func (h *Handler) UsersTable(w http.ResponseWriter, r *http.Request) {
	client, err := h.getClient()
	if err != nil || client == nil {
		h.renderPartialError(w, "Failed to load settings.")
		return
	}

	users, apiErr := client.ListUsers()
	if apiErr != nil {
		h.renderPartialError(w, apiErr.Error())
		return
	}

	h.render(w, "users-content.html", map[string]interface{}{
		"Title":      "Users",
		"ActivePage": "users",
		"Users":      users,
	})
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		h.renderToast(w, "Username is required.", "error")
		return
	}

	client, err := h.getClient()
	if err != nil || client == nil {
		h.renderToast(w, "Failed to load settings.", "error")
		return
	}

	_, apiErr := client.CreateUser(name, r.FormValue("displayName"), r.FormValue("email"), "https://robohash.org/"+name)
	if apiErr != nil {
		h.renderToast(w, apiErr.Error(), "error")
		return
	}

	h.renderToast(w, "User '"+name+"' created successfully!", "success")
}

func (h *Handler) RenameUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	oldID := r.FormValue("oldId")
	newName := r.FormValue("newName")

	if oldID == "" || newName == "" {
		h.renderToast(w, "User ID and new name are required.", "error")
		return
	}

	client, err := h.getClient()
	if err != nil || client == nil {
		h.renderToast(w, "Failed to load settings.", "error")
		return
	}

	if _, apiErr := client.RenameUser(oldID, newName); apiErr != nil {
		h.renderToast(w, apiErr.Error(), "error")
		return
	}

	h.renderToast(w, "User renamed to '"+newName+"' successfully!", "success")
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		h.renderToast(w, "User ID is required.", "error")
		return
	}

	client, err := h.getClient()
	if err != nil || client == nil {
		h.renderToast(w, "Failed to load settings.", "error")
		return
	}

	if apiErr := client.DeleteUser(id); apiErr != nil {
		h.renderToast(w, apiErr.Error(), "error")
		return
	}

	h.renderToast(w, "User deleted successfully!", "success")
}
