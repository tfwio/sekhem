// +build !session

package config

import (
	"github.com/gin-gonic/gin"
)

func (c *Configuration) initServerLogin(router *gin.Engine) bool {
	// fmt.Println("--> LOGON SESSIONS NOT SUPPORTED")
	// do nothing
	return false
}
