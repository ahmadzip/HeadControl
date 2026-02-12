package model

type Settings struct {
	ID        int    `json:"id"`
	BaseURL   string `json:"base_url"`
	APIKey    string `json:"api_key"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type User struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	CreatedAt     string `json:"createdAt"`
	DisplayName   string `json:"displayName"`
	Email         string `json:"email"`
	ProfilePicURL string `json:"profilePicUrl"`
}

type Node struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	GivenName       string   `json:"givenName"`
	IPAddresses     []string `json:"ipAddresses"`
	User            *User    `json:"user"`
	LastSeen        string   `json:"lastSeen"`
	Expiry          string   `json:"expiry"`
	CreatedAt       string   `json:"createdAt"`
	RegisterMethod  string   `json:"registerMethod"`
	Online          bool     `json:"online"`
	ApprovedRoutes  []string `json:"approvedRoutes"`
	AvailableRoutes []string `json:"availableRoutes"`
	Tags            []string `json:"tags"`
}

type DashboardStats struct {
	UserCount    int
	NodeCount    int
	OnlineNodes  int
	ExpiringSoon int
}

type UsersResponse struct {
	Users []User `json:"users"`
}

type UserResponse struct {
	User User `json:"user"`
}

type NodesResponse struct {
	Nodes []Node `json:"nodes"`
}

type NodeResponse struct {
	Node Node `json:"node"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
