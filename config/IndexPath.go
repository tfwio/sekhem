package config

// IndexPath — same as StaticPath, however we can call on something like
// `[target].json` to calculate/generate a file-index listing.
type IndexPath struct {
	Alias       string   `json:"alias,omitempty"` // this isn't supported really, but the idea is to use this as the name of our target as opposed to our root-directory name, or default server path (in Server.Path).
	Source      string   `json:"src"`
	Target      string   `json:"tgt"`
	Browsable   bool     `json:"nav"` // show directory file-listing in browser
	Servable    bool     `json:"serve"`
	IgnorePaths []string `json:"ignorePaths"` // absolute paths to ignore
	Extensions  []string `json:"spec"`        // file extensions to recognize; I.E.: the `Configuration.Extensions` .Name.
	path        string   // path as used in memory; we'll probably just ignore this guy.
}
