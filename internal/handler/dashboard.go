package handler

import (
	"headcontrol/internal/model"
	"net/http"
	"sort"
	"time"
)

func (h *Handler) DashboardPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	stats, recent, err := h.fetchDashboardData()
	if err != "" {
		h.renderPageWithError(w, r, "Dashboard", "dashboard", err)
		return
	}

	h.renderPage(w, r, "dashboard", map[string]interface{}{
		"Title":       "Dashboard",
		"ActivePage":  "dashboard",
		"Stats":       stats,
		"RecentNodes": recent,
	})
}

func (h *Handler) DashboardSummary(w http.ResponseWriter, r *http.Request) {
	stats, recent, err := h.fetchDashboardData()
	if err != "" {
		h.renderPartialError(w, err)
		return
	}

	h.render(w, "dashboard-content.html", map[string]interface{}{
		"Title":       "Dashboard",
		"ActivePage":  "dashboard",
		"Stats":       stats,
		"RecentNodes": recent,
	})
}

func (h *Handler) fetchDashboardData() (model.DashboardStats, []model.Node, string) {
	client, clientErr := h.getClient()
	if clientErr != nil || client == nil {
		return model.DashboardStats{}, nil, "Failed to load settings."
	}

	users, usersErr := client.ListUsers()
	if usersErr != nil {
		return model.DashboardStats{}, nil, usersErr.Error()
	}

	nodes, nodesErr := client.ListNodes()
	if nodesErr != nil {
		return model.DashboardStats{}, nil, nodesErr.Error()
	}

	online, expiring := 0, 0
	now := time.Now()
	for _, n := range nodes {
		if n.Online {
			online++
		}
		if t, err := time.Parse(time.RFC3339, n.Expiry); err == nil {
			if !t.IsZero() && t.After(now) && t.Before(now.Add(7*24*time.Hour)) {
				expiring++
			}
		}
	}

	stats := model.DashboardStats{
		UserCount:    len(users),
		NodeCount:    len(nodes),
		OnlineNodes:  online,
		ExpiringSoon: expiring,
	}

	sorted := make([]model.Node, len(nodes))
	copy(sorted, nodes)
	sort.Slice(sorted, func(i, j int) bool {
		ti, _ := time.Parse(time.RFC3339, sorted[i].LastSeen)
		tj, _ := time.Parse(time.RFC3339, sorted[j].LastSeen)
		return ti.After(tj)
	})
	if len(sorted) > 5 {
		sorted = sorted[:5]
	}

	return stats, sorted, ""
}
