package fsindex

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// // Hello uses import rsc.io/quote
// func Hello() string {
// 	return quote.hello()
// }

// sha1string func
func sha1String(pStrData string) string {
	hasher := sha1.New()
	hasher.Write([]byte(pStrData))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

// FileEntry ...
type FileEntry struct {
	Parent   *PathEntry // Parent directory
	FullPath string     // Complete directory path
	SHA1     string
}

// Abs ...Get the absolute path of a given directory.
func (refFileEntry *FileEntry) Abs() string {
	result, _ := filepath.Abs(refFileEntry.FullPath)
	return result
}

// Base ..
func (refFileEntry *FileEntry) Base() string {
	return filepath.Base(refFileEntry.FullPath)
}

// getSHA1 stores SHA1 hash on FileEntry and returns the result.
func (refFileEntry *FileEntry) getSHA1() string {
	refFileEntry.SHA1 = sha1String(refFileEntry.FullPath)
	return refFileEntry.SHA1
}

// Rooted returns the FileEntry.FullPath excluding the full root-path with exception to
// the root-directory name.
func (refFileEntry *FileEntry) Rooted(pRoot *PathEntry) string {
	return strings.Replace(refFileEntry.Abs(), pRoot.Abs(), pRoot.Base(), -1)
}

// PathSpec has to have a comment so there it is.
//
// This structure is basis for file/directory navigation wrapping folder/file structure in memory.
type PathSpec struct {
	FileEntry

	// Indicates a top-level directory
	isRoot bool

	// Child items
	Paths []PathEntry
	Files []FileEntry
}

// PathEntry ...
type PathEntry struct {
	PathSpec

	FileFilter  []FileSpec
	index       int32
	IgnorePaths []string
}

// IsIgnore ...
func (refPathEntry *PathEntry) IsIgnore(rootPathEntry *PathEntry) bool {
	for _, mNode := range rootPathEntry.IgnorePaths {
		if mNode == refPathEntry.Abs() {
			return true
		}
	}
	return false
}

// Refresh refreshes child directories and files.
func (refPathEntry *PathEntry) Refresh(rootPathEntry *PathEntry, counter *int32) {

	// create a reference node pointing to the tree-root
	var mRoot *PathEntry

	if refPathEntry.isRoot {
		mRoot = refPathEntry // rootPathEntry is `nil`
		// Absolute path for strings.Replace(…) functionality.
		for i := 0; i < len(refPathEntry.IgnorePaths); i++ {
			refPathEntry.IgnorePaths[i], _ = filepath.Abs(refPathEntry.IgnorePaths[i])
		}
	} else {
		mRoot = rootPathEntry // assign mRoot
	}

	refPathEntry.index = *counter // Assign index
	*counter++

	sha1 := refPathEntry.getSHA1()
	fmt.Printf("- %-4d: %s: %s\n", refPathEntry.index, sha1, refPathEntry.Rooted(mRoot))

	mPaths, mError := filepath.Glob(fmt.Sprintf("%s/*", refPathEntry.FullPath))

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
						Parent:   refPathEntry,
						FullPath: mFullPath,
						// SHA1: '',
					}
					refPathEntry.Files = append(refPathEntry.Files, child)
					println(fmt.Sprintf("  - %s", child.Base()))
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
						Parent:   refPathEntry,
						FullPath: mFullPath,
					},
					isRoot: false,
				},
			}

			if child.IsIgnore(mRoot) {

				fmt.Printf("- ignored: %s\n", child.FullPath)

			} else {

				child.Refresh(mRoot, counter)

				refPathEntry.Paths = append(refPathEntry.Paths, child)
			}

		}
	}
}

/////////////////////////////////////////////////////////////////////////////
// FileSpec
/////////////////////////////////////////////////////////////////////////////

// FileSpec structure.
type FileSpec struct {
	Name       string
	Extensions []string
}

// Match checks to see if an input file extention matches
// any of the file extensions defined in a given FileSpec.
func (refFileSpec *FileSpec) Match(input string) bool {

	fext := strings.ToLower(filepath.Ext(input))

	for i := 0; i < len(refFileSpec.Extensions); i++ {

		if refFileSpec.Extensions[i] == fext {
			return true
		}
	}
	return false
}

// var (
// 	xCounter   int32
// 	xpCounter  *int32
// 	localMedia = FileSpec{
// 		Name: "Media (images)",
// 		Extensions: []string{
// 			".bmp",
// 			".jpg",
// 			".svg",
// 			".png",
// 			".gif",
// 		},
// 	}
// 	localMarkdown = FileSpec{
// 		Name: "Markdown (hyper-text)",
// 		Extensions: []string{
// 			".md",
// 			".mmd",
// 		},
// 	}
// )

//func main() {
//	// flag.Parse()
//	// appRootPath := filepath.Dir(flag.Arg(0))
//	// root, _ := filepath.Abs(appRootPath)
//	// fmt.Println(fmt.Sprintf("root path: %s", root))
//	xCounter = 0
//	xpCounter = &xCounter
//	*xpCounter++
//	fmt.Printf("counter %d\n", *xpCounter)
//	*xpCounter++
//	fmt.Printf("counter %d\n", *xpCounter)
//	*xpCounter++
//	fmt.Printf("counter %d\n", *xpCounter)
//	startPath := "c:\\users\\tfwro\\.mmd"
//	rootPath := PathEntry{
//		PathSpec: PathSpec{
//			FileEntry: FileEntry{
//				Parent:   nil,
//				FullPath: startPath,
//			},
//			isRoot: true,
//		},
//		FileFilter: []FileSpec{localMedia, localMarkdown},
//		IgnorePaths: []string{
//			"c:\\users\\tfwro\\.mmd\\reveal.js",
//			"c:\\users\\tfwro\\.mmd\\.git",
//			"c:\\users\\tfwro\\.mmd\\.vscode",
//		},
//	}
//	xCounter = 0
//	rootPath.Refresh(nil, &xCounter)
//}
//
