package fsindex

// Default Settings
var (
	DefaultSettings = Settings{
		OmitRootNameFromPath: false,
	}
)

// Settings will slightly alter how the `Refresh` method runs.
type Settings struct {
	OmitRootNameFromPath bool
}

// Model is the same as PathEntry but with Settings
type Model struct {
	PathEntry
	SimpleModel
	Settings
}
