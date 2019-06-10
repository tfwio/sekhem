package fsindex

import (
	"path/filepath"
	"strings"
)

/////////////////////////////////////////////////////////////////////////////
// FileSpec
/////////////////////////////////////////////////////////////////////////////

// FileSpec structure.
type FileSpec struct {
	Name       string   `json:"name"`
	Extensions []string `json:"ext"`
}

// Match checks to see if an input file extention matches
// any of the file extensions defined in a given FileSpec.
func (f *FileSpec) Match(input string) bool {

	fext := strings.ToLower(filepath.Ext(input))

	for i := 0; i < len(f.Extensions); i++ {

		if f.Extensions[i] == fext {
			return true
		}
	}
	return false
}
