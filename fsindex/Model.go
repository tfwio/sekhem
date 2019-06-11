package fsindex

// Default Settings
var (
	Default = Settings{
		OmitRootNameFromPath: false,
	}
	currentSettings = Default
)

// Settings will slightly alter how the `Refresh` method runs.
type Settings struct {
	OmitRootNameFromPath bool
}

// Model is the same as PathEntry but with Settings
type Model struct {
	PathEntry
	Settings Settings
}
