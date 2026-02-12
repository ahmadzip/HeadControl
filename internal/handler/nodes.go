package handler

import (
	"net/http"
	"strings"
)

func (h *Handler) NodesPage(w http.ResponseWriter, r *http.Request) {
	client, err := h.getClient()
	if err != nil || client == nil {
		h.renderPageWithError(w, r, "Nodes", "nodes", "Failed to load settings.")
		return
	}

	nodes, apiErr := client.ListNodes()
	if apiErr != nil {
		h.renderPageWithError(w, r, "Nodes", "nodes", apiErr.Error())
		return
	}

	h.renderPage(w, r, "nodes", map[string]interface{}{
		"Title":      "Nodes",
		"ActivePage": "nodes",
		"Nodes":      nodes,
	})
}

func (h *Handler) NodesTable(w http.ResponseWriter, r *http.Request) {
	client, err := h.getClient()
	if err != nil || client == nil {
		h.renderPartialError(w, "Failed to load settings.")
		return
	}

	nodes, apiErr := client.ListNodes()
	if apiErr != nil {
		h.renderPartialError(w, apiErr.Error())
		return
	}

	h.render(w, "nodes-content.html", map[string]interface{}{
		"Title":      "Nodes",
		"ActivePage": "nodes",
		"Nodes":      nodes,
	})
}

func (h *Handler) NodeDetail(w http.ResponseWriter, r *http.Request) {
	nodeID := r.URL.Query().Get("id")
	if nodeID == "" {
		h.renderPartialError(w, "Node ID is required.")
		return
	}

	client, err := h.getClient()
	if err != nil || client == nil {
		h.renderPartialError(w, "Failed to load settings.")
		return
	}

	node, apiErr := client.GetNode(nodeID)
	if apiErr != nil {
		h.renderPartialError(w, apiErr.Error())
		return
	}

	h.render(w, "node-detail.html", map[string]interface{}{"Node": node})
}

func (h *Handler) RenameNode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	nodeID := r.FormValue("nodeId")
	newName := r.FormValue("newName")

	if nodeID == "" || newName == "" {
		h.renderToast(w, "Node ID and new name are required.", "error")
		return
	}

	client, err := h.getClient()
	if err != nil || client == nil {
		h.renderToast(w, "Failed to load settings.", "error")
		return
	}

	if _, apiErr := client.RenameNode(nodeID, newName); apiErr != nil {
		h.renderToast(w, apiErr.Error(), "error")
		return
	}

	h.renderToast(w, "Node renamed to '"+newName+"' successfully!", "success")
}

func (h *Handler) ExpireNode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	nodeID := r.FormValue("nodeId")
	if nodeID == "" {
		h.renderToast(w, "Node ID is required.", "error")
		return
	}

	client, err := h.getClient()
	if err != nil || client == nil {
		h.renderToast(w, "Failed to load settings.", "error")
		return
	}

	if _, apiErr := client.ExpireNode(nodeID); apiErr != nil {
		h.renderToast(w, apiErr.Error(), "error")
		return
	}

	h.renderToast(w, "Node expired successfully!", "success")
}

func (h *Handler) DeleteNode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	nodeID := r.FormValue("nodeId")
	if nodeID == "" {
		h.renderToast(w, "Node ID is required.", "error")
		return
	}

	client, err := h.getClient()
	if err != nil || client == nil {
		h.renderToast(w, "Failed to load settings.", "error")
		return
	}

	if apiErr := client.DeleteNode(nodeID); apiErr != nil {
		h.renderToast(w, apiErr.Error(), "error")
		return
	}

	h.renderToast(w, "Node deleted successfully!", "success")
}

func (h *Handler) SetNodeTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	nodeID := r.FormValue("nodeId")
	if nodeID == "" {
		h.renderToast(w, "Node ID is required.", "error")
		return
	}

	tags := splitCSV(r.FormValue("tags"))

	client, err := h.getClient()
	if err != nil || client == nil {
		h.renderToast(w, "Failed to load settings.", "error")
		return
	}

	if _, apiErr := client.SetNodeTags(nodeID, tags); apiErr != nil {
		h.renderToast(w, apiErr.Error(), "error")
		return
	}

	h.renderToast(w, "Tags updated successfully!", "success")
}

func (h *Handler) SetNodeRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	nodeID := r.FormValue("nodeId")
	if nodeID == "" {
		h.renderToast(w, "Node ID is required.", "error")
		return
	}

	routes := splitCSV(r.FormValue("routes"))

	client, err := h.getClient()
	if err != nil || client == nil {
		h.renderToast(w, "Failed to load settings.", "error")
		return
	}

	if _, apiErr := client.SetApprovedRoutes(nodeID, routes); apiErr != nil {
		h.renderToast(w, apiErr.Error(), "error")
		return
	}

	h.renderToast(w, "Routes approved successfully!", "success")
}

func splitCSV(raw string) []string {
	if raw == "" {
		return nil
	}
	var out []string
	for _, s := range strings.Split(raw, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}
