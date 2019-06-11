package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"tfw.io/Go/fsindex/fsindex"
	"tfw.io/Go/fsindex/util"
)

func createPathEntry(path string, spath string) fsindex.Model {
	println("- path for indexed files: ", configuration.GetPath(spath))
	// configure, createIndex, checkSimpleModel
	pe := fsindex.Model{
		PathEntry: fsindex.PathEntry{
			PathSpec: fsindex.PathSpec{
				FileEntry: fsindex.FileEntry{
					Parent:   nil,
					Name:     util.AbsBase(path),
					FullPath: util.Abs(path),
					SHA1:     util.Sha1String(path),
				},
				IsRoot: true},
			FauxPath:    configuration.GetPath(spath),
			FileFilter:  configuration.Extensions,
			IgnorePaths: []string{},
		},
		Settings: fsindex.Default,
	}
	return pe
}

func buildFileSystemModel(model *fsindex.Model, path string, spath string) {

	xCounter, fCounter = 0, 0

	mdl := fsindex.SimpleModel{}
	mdl.Create()

	handler := fsindex.Handlers{
		ChildPath: func(root *fsindex.Model, child *fsindex.PathEntry) bool {
			mdl.AddPath(root, child)
			return false
		},
		ChildFile: func(root *fsindex.Model, child *fsindex.FileEntry) bool {
			ext := strings.ToLower(filepath.Ext(child.FullPath))
			if ext == ".md" {
				// datestring := checkDateString(child.Base())
				mdl.AddFile(root, child)
			}
			fCounter++
			return false
		},
	}

	model.Refresh(model, &xCounter, &handler)

	// checkSimpleModel(&mdl)
}

func checkSimpleModel(mdl *fsindex.SimpleModel, pathEntry *fsindex.Model) {
	// map counters don't yield adequate length
	println("File map Count: ", len(mdl.File))
	println("Path map Count: ", len(mdl.Path))
	//
	println("File Count: ", fCounter)
	println("Path Count: ", xCounter)

	ref1 := &pathEntry.Paths[0].Files[0]
	println("some model: ", ref1.FullPath)
	println("parent:", ref1.Parent.FauxPath)
	fmt.Printf("looking in \"%s\" for files...\n", ref1.Parent.Base())
	for _, x := range ref1.Parent.Files {
		println("  -->", x.Path)
	}
}
