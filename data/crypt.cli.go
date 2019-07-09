package main

// tickers would help mitigate sessions
// https://gobyexample.com/tickers

// +build ignore

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"tfw.io/Go/fsindex/fsindex/ormus"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"tfw.io/Go/fsindex/util"
)

var (
	fdb      = flag.String("db", "data/ormus.db", "specify a database to use.")
	fSaltLen = flag.Int("s", 32, "provide default salt length")
	//
	fList = flag.NewFlagSet("list", flag.ExitOnError)
	//
	fCreate = flag.NewFlagSet("create", flag.ExitOnError)
	fcUser  = fCreate.String("u", "admin", "speficy username.")
	fcPass  = fCreate.String("p", "", "for validation and creation of a login profile.")
	fcSess  = fCreate.Bool("sess", false, "Create a session for the user.")
	//
	fValidate = flag.NewFlagSet("validate", flag.ExitOnError)
	fvUser    = fValidate.String("u", "admin", "speficy username.")
	fvPass    = fValidate.String("p", "", "for validation and creation of a login profile.")
	//
	//fvalid   = flag.String("V", "", "for validation and creation of a login profile.")
	//fsalt    = flag.String("salt", "", "[optional] supply salt and hash to validate -V <pass> (or fallback to db).")
	//fhash    = flag.String("hash", "", "[optional] supply salt and hash to validate -V <pass> (or fallback to db).")
)

func testDatabase() {
	db, err := gorm.Open("sqlite3", *fdb)
	defer db.Close()
	if err != nil {
		fmt.Printf("error loading empty database: %v\n", db)
	} else {
		println("success: opened empty database.")
	}
}

func main() {

	if len(os.Args) == 1 {
		flag.PrintDefaults()
		println()
		fmt.Printf("%s create (subcommand) args:\n", util.AbsBase(os.Args[0]))
		println()
		fCreate.PrintDefaults()
		println()
		fmt.Printf("%s validate (subcommand) args:\n", util.AbsBase(os.Args[0]))
		println()
		fValidate.PrintDefaults()
		return
	}

	ormus.SetSource(*fdb)
	ormus.EnsureTableUsers()
	ormus.EnsureTableSessions()

	switch strings.ToLower(os.Args[1]) {
	case "create":
		fCreate.Parse(os.Args[2:])
		if len(*fcPass) > 4 && len(*fcUser) > 3 {
			u := ormus.User{Name: *fcUser}
			if *fcSess {
				println("- Sesssion generation requested")
			}
			u.Create(*fcPass, *fSaltLen)
			println("- Sesssion generation requested")
			u.CreateSession()
		} else {
			fmt.Printf("- username %s; pass %s\n", *fcUser, *fcPass)
			println("- username must be > len(3) chars long")
			println("- you must supply a password > len(4) chars long")
		}
	case "validate":
		fValidate.Parse(os.Args[2:])
		if len(*fvPass) > 4 && len(*fvUser) > 3 {
			u := ormus.User{Name: *fvUser}
			result := u.Validate(*fvPass)
			fmt.Printf("Result: %v \n", result)
		} else {
			fmt.Printf("- username %s; pass %s\n", *fvUser, *fvPass)
			println("- username must be > len(3) chars long")
			println("- you must supply a password > len(4) chars long")
		}
	case "list":
		var users []ormus.User
		var sessions []ormus.Session
		db, err := gorm.Open("sqlite3", *fdb)
		defer db.Close()
		if err != nil {
			fmt.Printf("error loading database: %v\n", db)
		} else {
			// list users
			// db.Where("[name] like ?").Find(&users)
			db.Find(&users)
			fmt.Printf("- found %d entries\n", len(users))
			usermap := make(map[int64]*ormus.User)
			for m, x := range users {
				fmt.Printf("  - %04d: %s\n", m, x.Name)
				fmt.Printf("    % 4d: %v\n", ' ', x)
				usermap[x.ID] = &x
			}
			// list sessions
			db.Find(&sessions)
			fmt.Printf("- found %d entries\n", len(sessions))
			for _, x := range sessions {
				// fmt.Printf("  - %04d: %s\n", m, x.Name)
				fmt.Printf("--> \"%16s\" %s, %s\n", usermap[x.UserID].Name, x.Created.Format("20060102_1504.05"), x.SessID)
			}
		}
	default:
		println()
		flag.PrintDefaults()
	}

}
