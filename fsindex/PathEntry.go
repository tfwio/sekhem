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

// Top gets the top-most, root path entry.
func (p *PathEntry) Top() *PathEntry {
	mRef := p
	for !mRef.IsRoot {
		mRef = p
	}
	return mRef
}

// IsIgnore ...
func (p *PathEntry) IsIgnore(r *PathEntry) bool {
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
func (p *PathEntry) Review(cbPath *CBPath, cbFile *CBFile) {
	for _, child := range p.Paths {
		if cbPath != nil {
			if (*cbPath)(child.Parent, &child) {
				return
			}
			child.Review(cbPath, cbFile)
		}
	}
	for _, child := range p.Files {
		if cbPath != nil {
			if (*cbFile)(child.Parent, &child) {
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
func (p *PathEntry) Refresh(rootPathEntry *PathEntry, counter *(int32), cbPath *CBPath, cbFile *CBFile) {

	// create a reference node pointing to the tree-root
	var mRoot *PathEntry

	// if the first parent element is root, we need to build some
	// reference memory (dictionary of ignore-paths).
	if p.IsRoot {
		mRoot = p
		// build absolute path list to ignore.
		for i := 0; i < len(p.IgnorePaths); i++ {
			p.IgnorePaths[i], _ = filepath.Abs(p.IgnorePaths[i])
		}
	} else {
		mRoot = rootPathEntry // assign mRoot
	}

	p.Index = *counter // Assign index
	*counter++

	if cbPath != nil {
		if (*cbPath)(mRoot, p) {
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

			for i := 0; i < len(mRoot.FileFilter); i++ {

				if mRoot.FileFilter[i].Match(mFullPath) {

					var child = FileEntry{
						Parent:   p,
						FullPath: mFullPath,
						SHA1:     util.Sha1String(mFullPath),
					}
					// Rooted only works once .FullPath is set.
					child.Path = util.UnixSlash(util.Cat(mRoot.FauxPath, "/", child.Rooted(p)))
					p.Files = append(p.Files, child)
					if cbFile != nil {
						if (*cbFile)(mRoot, &child) {
							return
						}
						// println(fmt.Sprintf("  - %s", child.Base()))
					}

				}
			}

			// if !isMediaExclude(mPath.Name()) {
			// indexMediaModelPaths(mPathAbs)
			// }
			// } else {
			// 	mFileAbs, _ := filepath.Abs(filepath.Join(mAbs, mPath.Name()))
			// 	if isMediaFile(mFileAbs) {
			// 		MediaFiles = append(MediaFiles, mFileAbs)
			// 	}
			// }
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
						Parent:   p,
						FullPath: mFullPath,
						SHA1:     util.Sha1String(mFullPath),
					},
					IsRoot: false,
				},
			}
			child.Path = util.UnixSlash(util.Cat(mRoot.FauxPath, "/", child.Rooted(p)))

			if child.IsIgnore(mRoot) {

				fmt.Printf("- ignored: %s\n", child.FullPath)

			} else {

				child.Refresh(mRoot, counter, cbPath, cbFile)

				p.Paths = append(p.Paths, child)
			}

		}
	}
}

// Refresh1 refreshes child directories and files.
// parameter `rootPathEntry`: root-path entry.
// parameter `counter (*int32)`: pointer to our indexing integer (counter).
// parameter `callback (RefreshAction)` is a method (if defined) which
//                                      can be used arbitrarily.
func (p *PathEntry) Refresh1(rootPathEntry *PathEntry, counter *(int32), handler *Handlers) {

	// create a reference node pointing to the tree-root
	var mRoot *PathEntry

	// if the first parent element is root, we need to build some
	// reference memory (dictionary of ignore-paths).
	if p.IsRoot {
		mRoot = p // rootPathEntry is `nil`
		// Absolute path for strings.Replace(…) functionality.
		for i := 0; i < len(p.IgnorePaths); i++ {
			p.IgnorePaths[i], _ = filepath.Abs(p.IgnorePaths[i])
		}
	} else {
		mRoot = rootPathEntry // assign mRoot
	}

	p.Index = *counter // Assign index
	*counter++

	if handler != nil {
		if handler.ChildPath(mRoot, p) {
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

			for i := 0; i < len(mRoot.FileFilter); i++ {

				if mRoot.FileFilter[i].Match(mFullPath) {

					var child = FileEntry{
						Parent:   p,
						FullPath: mFullPath,
						SHA1:     util.Sha1String(mFullPath),
						Name:     filepath.Base(mFullPath),
					}
					child.Path = util.UnixSlash(util.Cat(mRoot.FauxPath, "/", child.Rooted(p)))
					p.Files = append(p.Files, child)
					if handler != nil {
						if handler.ChildFile(mRoot, &child) {
							return
						}
						// println(fmt.Sprintf("  - %s", child.Base()))
					}

				}
			}

			// if !isMediaExclude(mPath.Name()) {
			// indexMediaModelPaths(mPathAbs)
			// }
			// } else {
			// 	mFileAbs, _ := filepath.Abs(filepath.Join(mAbs, mPath.Name()))
			// 	if isMediaFile(mFileAbs) {
			// 		MediaFiles = append(MediaFiles, mFileAbs)
			// 	}
			// }
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
						Parent:   p,
						FullPath: mFullPath,
						SHA1:     util.Sha1String(mFullPath),
						Name:     filepath.Base(mFullPath),
					},
					IsRoot: false,
				},
			}
			child.Path = util.UnixSlash(util.Cat(mRoot.FauxPath, "/", child.Rooted(p)))

			if child.IsIgnore(mRoot) {

				fmt.Printf("- ignored: %s\n", child.FullPath)

			} else {

				child.Refresh1(mRoot, counter, handler)

				p.Paths = append(p.Paths, child)
			}

		}
	}
}
