package config

// StaticPath is a definition for directories we'll
// allow into the app, preferably by way of JSON config.
type StaticPath struct {
	Source string `json:"src"`
	Target string `json:"tgt"`
	// show directory file-listing in browser
	Browsable bool `json:"nav"`
}
