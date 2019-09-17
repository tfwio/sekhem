package config

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tfwio/srv/util"
)

// JSONIndex — a simple container for JSON.
type JSONIndex struct {
	Index []string `json:"index"`
}

func (c *Configuration) serveJSONIndex(g *gin.Context) {

	loggedIn := false
	if loggedIn1, success := g.Get("valid"); success {
		loggedIn = loggedIn1.(bool)
	}

	xdata := JSONIndex{} // xdata indexes is just a string array map.
	xdata.Index = []string{}
	for _, path := range c.Indexes {
		// fmt.Printf("--> requires-login(%v) and logged-in(%v)\n", path.RequiresLogin, loggedIn)
		if path.RequiresLogin && loggedIn {
			xdata.Index = append(xdata.Index, util.WReap("/", "json", util.AbsBase(path.Source)))
		} else if !path.RequiresLogin {
			xdata.Index = append(xdata.Index, util.WReap("/", "json", util.AbsBase(path.Source)))
		}
	}
	g.JSON(http.StatusOK, xdata)
}

func (c *Configuration) serveModelIndex(router *gin.Engine) {
	println("location indexes #2: primary")
	for _, path := range c.Indexes {
		jsonpath := util.WReap("/", "json", util.AbsBase(path.Source))
		modelpath := util.WReap("/", path.Target)
		fmt.Printf("  > Target = %-18s, json = %s,  Source = %s\n", modelpath, c.GetPath(jsonpath), path.Source)
		modelpath = c.getIndexTarget(&path)

		if path.Servable {
			router.StaticFS(modelpath, gin.Dir(util.Abs(path.Source), path.Browsable))
		}
	}
	router.Any("/json/:route", c.serveJSON)
	router.Any("/refresh/:route", c.refreshRouteJSON)
	router.Any("/tag/:route/*action", func(g *gin.Context) { TagHandler(c, g) })
	router.Any("/jtag/:route/*action", func(g *gin.Context) { TagHandlerJSON(c, g) })
}

func (c *Configuration) serveJSON(ctx *gin.Context) {

	mroute := ctx.Param("route")

	if c.hasModel(mroute) {
		mmdl := mdlMap[mroute]
		ctx.JSON(http.StatusOK, &mmdl.PathEntry)
	} else {
		jsi := JSONIndex{Index: []string{fmt.Sprintf("COULD NOT find model for index: %s", mroute)}}
		ctx.JSON(http.StatusNotFound, &jsi)
		fmt.Printf("--> COULD NOT FIND ROUTE %s\n", mroute)
	}
}

func (c *Configuration) refreshRouteJSON(g *gin.Context) {
	mroute := g.Param("route")
	jsi := JSONIndex{Index: []string{fmt.Sprintf("FOUND model for index: %s", mroute)}}
	if ndx, ok := c.indexFromTarget(mroute), c.hasModel(mroute); ok && ndx != nil {
		c.initializeModel(ndx)
		g.JSON(http.StatusOK, jsi)
		return
	}
	jsi = JSONIndex{Index: []string{fmt.Sprintf("COULD NOT find model for index: %s", mroute)}}
	g.JSON(http.StatusOK, &jsi)
	fmt.Printf("ERROR> COULD NOT find model for index: %s\n", mroute)
}
