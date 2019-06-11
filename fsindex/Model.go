package fsindex

// PathType has a Refresh method
type PathType interface {
	Refresh(rootPathEntry *PathEntry, counter *(int32), handler *Handlers)
}

// Model is the same as PathEntry but with Settings
type Model struct {
	PathEntry
	Settings Settings
}
