package fsindex

// Default Settings
var (
	DefaultSettings = Settings{
		OmitRootNameFromPath:       false,
		StripFileExtensionFromName: true,
		StripFileExtensionFromPath: false,
		UnknownCharsToDash:         false,
	}
)

// Settings will slightly alter how the `Refresh` method runs.
type Settings struct {
	OmitRootNameFromPath       bool
	StripFileExtensionFromName bool
	StripFileExtensionFromPath bool
	UnknownCharsToDash         bool
}

// Model is the same as PathEntry but with Settings
type Model struct {
	PathEntry
	SimpleModel `json:"-"`
	Settings    `json:"-"`
}
