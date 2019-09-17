package fsindex

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tfwio/srv/util"
)

// FileEntry ...
type FileEntry struct {
	Parent    *PathEntry `json:"-"`              // Parent directory
	Name      string     `json:"name,omitempty"` //
	FullPath  string     `json:"-"`              // Complete directory path
	SHA1      string     `json:"sha1,omitempty"`
	Path      string     `json:"path,omitempty"`
	Extension string     `json:"ext,omitempty"`
	Mod       time.Time  `json:"mod"`
}

// Abs ...Get the absolute path of a given directory.
func (f *FileEntry) Abs() string {
	result, _ := filepath.Abs(f.FullPath)
	return result
}

// Base ..
func (f *FileEntry) Base() string {
	return filepath.Base(f.FullPath)
}

// GetSHA1 stores SHA1 hash on FileEntry and returns the result.
func (f *FileEntry) GetSHA1() string {
	f.SHA1 = util.Sha1String(f.FullPath)
	return f.SHA1
}

// Rooted returns the FileEntry.FullPath excluding the full root-path with exception to
// the root-directory name.  It applies `Settings.OmitRootNameFromPath`.
func (f *FileEntry) Rooted(r *Model) string {
	result := ""
	if r.Settings.OmitRootNameFromPath {
		result = strings.Replace(f.Abs(), r.Abs(), r.Base(), -1)
	} else {
		result = strings.Replace(f.Abs(), r.Abs(), r.Name, -1)
	}
	return strings.Trim(result, "/")
}

// RootedPath applies additional filtering on `FileEntry` such as
// `Settings.UnknownCharsToDash` and `Settings.OmitRootNameFromPath`.
func (f *FileEntry) RootedPath(r *Model) string {
	result := f.Rooted(r)
	if r.Settings.UnknownCharsToDash {
		result = util.Space2Dash(result)
	}
	return result
}

// Modified gets the file-system modified time.
func (f *FileEntry) Modified() string {
	mfile, err := os.Stat(f.FullPath)
	if err != nil {
		return err.Error()
	}
	return mfile.ModTime().Format("2006-01-02")
}
