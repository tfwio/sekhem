package config

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
	"github.com/tfwio/sekhem/fsindex/ormus"
	"github.com/tfwio/sekhem/util"
)

func loginCheckUP(user string, pass string) bool {
	result := false
	if len(user) > 4 && len(pass) > 3 {
		u := ormus.User{Name: user}
		result = u.Validate(pass)
		//fmt.Printf("Result: %v \n", result)
	} else {
		//fmt.Printf("- username %s; pass %s\n", user, pass)
		println("- username must be > len(3) chars long")
		println("- you must supply a password > len(4) chars long")
	}
	return result
}

// LogonIsValid checks to see if a user is logged in.
//
func (c *Configuration) LogonIsValid(g *gin.Context) bool {
	ormus.EnsureTableSessions()
	ormus.EnsureTableUsers()

	sessLookup := fmt.Sprintf("SKM%s", c.Port)
	sessID, err := g.Cookie(sessLookup)
	if err == nil {
		fmt.Println("Error loading cookie")
		return false
	}
	fmt.Printf("found session %s\n", sessID)
	dbPath := util.CatPath(util.GetDirectory(util.Abs(DefaultConfigFile)), c.Database)
	db, err := gorm.Open(c.DatabaseType, dbPath)

	defer db.Close()
	if err != nil {
		println("well we've opened the database.")
	} else {
		fmt.Printf("error loading database: %v\n", db)
		return false
	}

	saltString := util.NewSaltString(64)
	bv := fmt.Sprintf("%s", saltString) == "1"
	return bv
}
