package ormus

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/tfwio/sekhem/util"
)

var (
	datasource string
)

// _ "github.com/mattn/go-sqlite3"

// User structure
type User struct {
	ID   int64  `gorm:"auto_increment;unique_index;primary_key;column:id"`
	Name string `gorm:"size:27;column:user"`
	Salt string `gorm:"size:432;column:salt"`
	Hash string `gorm:"size:432;column:hash"`
}

// SetSource allows a external library to set the local datasource.
func SetSource(source string) {
	datasource = source
}

// EnsureTableUsers creates table [users] if not exist.
func EnsureTableUsers() {
	db, err := gorm.Open("sqlite3", datasource)
	defer db.Close()
	if err != nil {
		fmt.Printf("error: ensuring database: users\n")
	} else {
		var u User
		if !db.HasTable(u) {
			db.CreateTable(u)
		}
	}

}

// TableName Set User's table name to be `users`
func (User) TableName() string {
	return "users"
}

/* http://jinzhu.me/gorm/crud.html#query */

// Get a user by name.
// there is also a non User-bound version of this named `ormus.GetUser(string)`.
// With this method, we just supply a name to the user structure.
func (u *User) Get() {
	db, err := gorm.Open("sqlite3", datasource)
	defer db.Close()
	if err != nil {
		fmt.Printf("- error getting user\n")
	} else {
		db.Where("name=?").First(u)
		fmt.Printf("- success getting user\n")
	}
}

// CreateSession is a test to attempt to save a session into the sessions table.
func (u *User) CreateSession() {
	sess := Session{Created: time.Now(), UserID: u.ID, SessID: util.NewSaltString(32)}
	db, err := gorm.Open("sqlite3", datasource)
	defer db.Close()
	if err != nil {
		println("unexpected error creating ")
	} else {
		db.Create(&sess)
		fmt.Printf("- session should have been created; sessid=%s\n", sess.SessID)
	}
}

// GetSessions gets sessions existing sessions for a user
func (u *User) GetSessions() {
	sess := Session{Created: time.Now(), UserID: u.ID, SessID: util.NewSaltString(32)}
	db, err := gorm.Open("sqlite3", datasource)
	defer db.Close()
	if err != nil {
		println("unexpected error creating ")
	} else {
		db.Create(&sess)
		fmt.Printf("- session should have been created; sessid=%s\n", sess.SessID)
	}
}

// Create http://jinzhu.me/gorm/crud.html#query
func (u *User) Create(pass string, saltSize int) {

	bsalt := util.NewSaltCSRNG(saltSize)
	u.Salt = util.BytesToBase64(bsalt)
	u.Hash = util.BytesToBase64(util.GetPasswordHash(pass, bsalt))

	fmt.Printf("Name: %s\nsalt: %s\nhash: %s\n", u.Name, u.Salt, u.Hash)

	db, err := gorm.Open("sqlite3", datasource) // defer db.Close()
	if err != nil {
		fmt.Printf("error: loading database: %v\n", db)
	} else {
		println("success: loading database.")
	}

	if !db.HasTable(u) {
		db.CreateTable(u)
	}

	var tempUser User // check if the record exists
	db.Where("user=?", u.Name).First(&tempUser)
	defer db.Close()
	if tempUser.Name != "" {
		fmt.Printf("- User exists! (user name taken)\n")
	} else {
		db.Create(u)
		fmt.Printf("- User created!\n")
	}
}

// validate checks against a provided salt and hash.
// This method does not actually look anything up from a database.
//
// Salt and Hash MUST BE PRESENT before calling!
func (u *User) validate(pass string) bool {
	result := util.CheckPassword(
		pass,
		util.FromBase64(u.Salt),
		util.FromBase64(u.Hash))
	return result
}

// GetUser by name.
func GetUser(name string) (*gorm.DB, User) {
	// open the database
	var tempUser User
	db, err := gorm.Open("sqlite3", datasource) //
	if err != nil {
		fmt.Printf("error: loading database: %v\n", db)
		db.Close()
		return nil, User{}
	}
	db.Where("user=?", name).First(&tempUser)
	return db, tempUser
}

// Validate checks against a provided salt and hash.
func (u *User) Validate(pass string) bool {
	// open the database
	db, err := gorm.Open("sqlite3", datasource) //
	if err != nil {
		fmt.Printf("error: loading database: %v\n", db)
		db.Close()
		return false
	}

	result := false
	var tempUser User // check if the record exists
	db.Where("user=?", u.Name).First(&tempUser)
	defer db.Close()
	if tempUser.Name != "" {
		result = tempUser.validate(pass)
	} else {
		println("- no user found.")
	}
	return result
}
