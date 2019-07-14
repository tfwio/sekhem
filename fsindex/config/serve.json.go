package config

// JSONIndex — a simple container for JSON.
type JSONIndex struct {
	Index []string `json:"index"`
}

// LogonModel responds to a login action such as "/login/" or (perhaps) "/login-refresh/"
type LogonModel struct {
	Action string `json:"action"`
	Status bool   `json:"status"`
	Detail string `json:"detail"`
}
