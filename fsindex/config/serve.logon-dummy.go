// +build !session

package config

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (c *Configuration) initServerLogin(router *gin.Engine) {
	fmt.Println("--> LOGON SESSIONS NOT SUPPORTED")
	router.Use(c.sessMiddleware)
	// do nothing
}

func (c *Configuration) sessMiddleware(g *gin.Context) {
	g.Set("valid", true) // pretend user is always logged in
}
