package fsindex

// CBPath is a simple callback;
// if you return true, then the caller function immediately returns.
type CBPath func(*PathEntry, *PathEntry) bool // (*interface{}, error)
// CBFile is a simple callback
// if you return true, then the caller function immediately returns.
type CBFile func(*PathEntry, *FileEntry) bool // (*interface{}, error)

// FileHandler is a simple callback.
type FileHandler func(*PathEntry, *FileEntry) bool

// PathHandler is a simple callback.
type PathHandler func(*PathEntry, *PathEntry) bool

// Handlers contains simple callbacks.
type Handlers struct {
	ChildPath PathHandler
	ChildFile FileHandler
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
