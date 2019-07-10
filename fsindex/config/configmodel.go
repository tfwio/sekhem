package config

/**
 * configmodel.go provides private methods for configuring `fsindex.Model` models
 * so that they can be served to JSON content.
 */

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tfwio/sekhem/fsindex"
	"github.com/tfwio/sekhem/util"
)

func buildFileSystemModel(model *fsindex.Model) {

	xCounter, fCounter = 0, 0

	model.SimpleModel = fsindex.SimpleModel{}
	model.CreateMaps()

	handler := fsindex.Handlers{
		ChildPath: func(root *fsindex.Model, child *fsindex.PathEntry) bool {
			model.AddPath(root, child)
			return false
		},
		ChildFile: func(root *fsindex.Model, child *fsindex.FileEntry) bool {
			ext := strings.ToLower(filepath.Ext(child.FullPath))
			if ext == ".md" {
				// datestring := checkDateString(child.Base())
				model.AddFile(root, child)
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

func (c *Configuration) getModelIndex(mdl *fsindex.Model) (int, bool) {
	for index, mMdl := range models {
		if mMdl.FullPath == mdl.FullPath {
			return index, true
		}
	}
	return -1, false
}

func (c *Configuration) getSimpleIndexTarget(path *IndexPath) string {
	return util.WReap("/", util.AbsBase(path.Source))
}

func (c *Configuration) getIndexTarget(path *IndexPath) string {
	modelpath := util.WReap("/", path.Target)
	if !c.IndexCfg.OmitRootNameFromPath {
		modelpath = util.WReap("/", path.Target, util.AbsBase(path.Source))
	}
	return modelpath
}

func (c *Configuration) hasModel(route string) bool {

	if _, ok := mdlMap[route]; ok {
		return true
	}
	return false
}

func (c *Configuration) indexFromTarget(route string) *IndexPath {
	inputTarget := util.WReap("/", route)
	for _, x := range c.Indexes {
		simpleIndexTarget := c.getSimpleIndexTarget(&x)
		if inputTarget == simpleIndexTarget {
			return &x
		}
	}
	return nil
}

func (c *Configuration) initializeModels() {
	for _, path := range c.Indexes {
		c.initializeModel(&path)
	}
}

func (c *Configuration) initializeModel(path *IndexPath) {

	fmt.Printf("--> indexing: %s\n", path.Target)
	model := c.createEntry(*path, c.IndexCfg)
	if _, ok := mdlMap[util.AbsBase(path.Source)]; !ok {
		models = append(models, model)
	} else {
		if index, ok := c.getModelIndex(&model); ok {
			models[index] = model
			println("Injecting memory-Model %s at index %d", mdlMap[model.Name].Name, index)
		}
	}
	if index, ok := c.getModelIndex(&model); ok {
		mdlMap[model.Name] = &models[index]
	} else {
		panic("Could not find memory-Model")
	}
}

func (c *Configuration) createEntry(path IndexPath, settings fsindex.Settings) fsindex.Model {

	// configure, createIndex, checkSimpleModel
	pe := fsindex.Model{
		PathEntry: fsindex.PathEntry{
			PathSpec: fsindex.PathSpec{
				FileEntry: fsindex.FileEntry{
					Parent:   nil,
					Name:     util.AbsBase(path.Source),
					FullPath: util.Abs(path.Source),
					SHA1:     util.Sha1String(path.Source),
				},
				IsRoot: true},
			FauxPath:    c.GetPath(path.Target),
			FileFilter:  c.Extensions,
			IgnorePaths: []string{},
		},
		Settings: settings,
	}
	buildFileSystemModel(&pe)
	return pe
}
