package config

// JSONIndex — a simple container for JSON.
type JSONIndex struct {
	Index []string `json:"index"`
}

// LogonModel responds to a login action such as "/login/" or (perhaps) "/login-refresh/"
type LogonModel struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}
