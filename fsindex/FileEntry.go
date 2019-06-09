package fsindex

import (
	"os"
	"path/filepath"
	"strings"

	"tfw.io/Go/fsindex/util"
)

// FileEntry ...
type FileEntry struct {
	Parent   *PathEntry `json:"-"`              // Parent directory
	Name     string     `json:"name,omitempty"` //
	FullPath string     `json:"-"`              // Complete directory path
	SHA1     string     `json:"sha1,omitempty"`
	Path     string     `json:"path,omitempty"`
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
// the root-directory name.
func (f *FileEntry) Rooted(r *PathEntry) string {
	return strings.Replace(f.Abs(), r.Abs(), r.Base(), -1)
}

// Modified gets the file-system modified time.
func (f *FileEntry) Modified() string {
	mfile, err := os.Stat(f.FullPath)
	if err != nil {
		return err.Error()
	}
	return mfile.ModTime().Format("2006-01-02")
}
