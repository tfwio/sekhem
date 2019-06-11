package fsindex

// SimpleModel collects our indexes
type SimpleModel struct {
	File     map[string]*FileEntry
	FileSHA1 map[string]*FileEntry
	Path     map[string]*PathEntry
	PathSHA1 map[string]*PathEntry
}

// Create makes a new data-set.
func (m *SimpleModel) Create() {
	m.File = make(map[string]*FileEntry)
	m.FileSHA1 = make(map[string]*FileEntry)
	m.Path = make(map[string]*PathEntry)
	m.PathSHA1 = make(map[string]*PathEntry)
}

// Reset destroys all top level items (if hierarchical) in the maps.
func (m *SimpleModel) Reset() {
	if m.File != nil {
		for k := range m.File {
			delete(m.File, k)
		}
	}
	if m.FileSHA1 != nil {
		for k := range m.FileSHA1 {
			delete(m.FileSHA1, k)
		}
	}
	if m.Path != nil {
		for k := range m.Path {
			delete(m.Path, k)
		}
	}
	if m.PathSHA1 != nil {
		for k := range m.PathSHA1 {
			delete(m.PathSHA1, k)
		}
	}
}

// AddPath is a callback per PathEntry.
// It adds each PathEntry to a flat (non-hierarchical) map (dictionary).
func (m *SimpleModel) AddPath(p *Model, c *PathEntry) {
	m.Path[c.Rooted(p)] = c
	m.PathSHA1[c.SHA1] = c
}

// AddFile is a callback per FileEntry.
// It adds each FileEntry to a flat (non-hierarchical) map (dictionary).
func (m *SimpleModel) AddFile(p *Model, c *FileEntry) {
	m.File[c.Rooted(p)] = c
	m.FileSHA1[c.SHA1] = c
}
