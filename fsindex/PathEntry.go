package fsindex

import (
	"fmt"
	"os"
	"path/filepath"

	"tfw.io/Go/fsindex/util"
)

// PathEntry ...
type PathEntry struct {
	PathSpec

	FileFilter  []FileSpec `json:"-"`
	Index       int32      `json:"id"`
	IgnorePaths []string   `json:"-"`
	// FauxPath is only set on the root item and is
	// used to portray a URI from a relative path.
	FauxPath string `json:"uri,omitempty"`
}

// Info prints out some PathEntry info, of course.
func (p *PathEntry) Info() {
	println("- check name:", p.Name)
	println("- check sha1:", p.SHA1)
	println("- check path:", p.FauxPath)
	// print(fmt.Sprintf("- JSON index: %s\n", util.AbsBase(util.Abs(p.Source))))
	for _, x := range p.FileFilter {
		print(fmt.Sprintf("  - got extension: %s\n", x.Name))
	}
}

// Top gets the top-most, root path entry.
func (p *PathEntry) Top() *PathEntry {
	mRef := p
	for !mRef.IsRoot {
		mRef = p
	}
	return mRef
}

// IsIgnore ...
func (p *PathEntry) IsIgnore(r *Model) bool {
	for _, mNode := range r.IgnorePaths {
		if mNode == p.Abs() {
			return true
		}
	}
	return false
}

// Review is similar to `Refresh()` except we don't rebuild the graph.
// here, we're just linking the callbacks, directories are listed before
// files like the `Refresh()` method.
func (p *PathEntry) Review(mRoot *Model, cbPath *CBPath, cbFile *CBFile) {
	for _, child := range p.Paths {
		if cbPath != nil {
			if (*cbPath)(mRoot, &child) {
				return
			}
			child.Review(mRoot, cbPath, cbFile)
		}
	}
	for _, child := range p.Files {
		if cbPath != nil {
			if (*cbFile)(mRoot, &child) {
				return
			}
		}
	}
}

// Refresh refreshes child directories and files.
// parameter `rootPathEntry`: root-path entry.
// parameter `counter (*int32)`: pointer to our indexing integer (counter).
// parameter `callback (RefreshAction)` is a method (if defined) which
//                                      can be used arbitrarily.
func (p *PathEntry) Refresh(model *Model, counter *(int32), handler *Handlers) {

	if p.IsRoot {
		for i := 0; i < len(p.IgnorePaths); i++ {
			p.IgnorePaths[i] = util.Abs(p.IgnorePaths[i])
		}
	}

	p.Index = *counter // Assign index
	*counter++

	if handler != nil {
		if handler.ChildPath(model, p) {
			return
		}
	}

	mPaths, mError := filepath.Glob(fmt.Sprintf("%s/*", p.FullPath))

	if mError != nil {
		fmt.Println("error in path:", mError)
		return
	}

	// FILE PATHS
	for _, mFullPath := range mPaths {

		fileinfo, err := os.Stat(mFullPath)
		if os.IsNotExist(err) {
			fmt.Println("Error reading file")
			return
		}

		if !fileinfo.IsDir() {

			for i := 0; i < len(model.FileFilter); i++ {

				if model.FileFilter[i].Match(mFullPath) {
					var child = FileEntry{
						Parent:    p,
						FullPath:  mFullPath,
						SHA1:      util.Sha1String(mFullPath),
						Name:      util.StripFileExtension(filepath.Base(mFullPath)),
						Extension: filepath.Ext(mFullPath),
						Mod:       fileinfo.ModTime(),
					}
					child.Path = util.UnixSlash(util.Cat(model.FauxPath, "/", child.Rooted(model)))
					p.Files = append(p.Files, child)
					if handler != nil {
						if handler.ChildFile(model, &child) {
							return
						}
					}

				}
			}

		}
	}

	// DIRECTORY PATHS
	for _, mFullPath := range mPaths {

		fileinfo, err := os.Stat(mFullPath)
		if os.IsNotExist(err) {
			fmt.Println("Error reading file")
			return
		}

		if fileinfo.IsDir() {

			var child = PathEntry{
				PathSpec: PathSpec{
					FileEntry: FileEntry{
						Parent:    p,
						FullPath:  mFullPath,
						SHA1:      util.Sha1String(mFullPath),
						Name:      util.StripFileExtension(filepath.Base(mFullPath)),
						Extension: filepath.Ext(mFullPath),
						Mod:       fileinfo.ModTime(),
					},
					IsRoot: false,
				},
			}
			child.Path = util.UnixSlash(util.Cat(model.FauxPath, "/", child.Rooted(model)))

			if !child.IsIgnore(model) {
				child.Refresh(model, counter, handler)
				p.Paths = append(p.Paths, child)
			}

		}
	}
}
