package fsindex

// Default Settings
var (
	DefaultSettings = Settings{
		OmitRootNameFromPath:       false,
		StripFileExtensionFromName: true,
		UnknownCharsToDash:         false,
		HardLinks:                  false,
	}
)

// Settings will slightly alter how the `Refresh` method runs.
// Eventually, we'll convert this to flag-usage in the console client.
type Settings struct {
	// OmitRootNameFromPath will strip the root directory-name from indexed path targets.
	// Only the default value of false is currently known to be working.
	// For example if true, a path converted to "http path": path-in: "c:/mypath/mysubdir/my-target-path", path-out: "/".
	// If set to (default) false: path-in: "c:/mypath/mysubdir/my-target-path", path-out: "/my-target-path".
	OmitRootNameFromPath       bool `json:"omit-root"`
	StripFileExtensionFromName bool `json:"strip-file-ext,omitempty"` // since default=true: "opmitempty".
	UnknownCharsToDash         bool `json:"space2dash,omitempty"`     // not uet supported.
	HardLinks                  bool `json:"hard-link"`                // this tells us weather or not to use full link-path such as `http://[server:port]/` when generating JSON.
}

// Model is the same as PathEntry but with Settings
type Model struct {
	PathEntry
	SimpleModel `json:"-"`
	Settings    `json:"-"`
}
