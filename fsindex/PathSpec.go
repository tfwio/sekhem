package fsindex

// PathSpec has to have a comment so there it is.
//
// This structure is basis for file/directory navigation wrapping folder/file structure in memory.
type PathSpec struct {
	FileEntry

	// Indicates a top-level directory
	IsRoot bool `json:"-"`

	// Child items
	Paths []PathEntry `json:"paths,omitempty"`
	Files []FileEntry `json:"files,omitempty"`
}
