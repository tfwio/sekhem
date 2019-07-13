package ormus

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tfwio/sekhem/util"
)

// User structure
type User struct {
	dbx
	ID   int64  `gorm:"auto_increment;unique_index;primary_key;column:id"`
	Name string `gorm:"size:27;column:user"`
	Salt string `gorm:"size:432;column:salt"`
	Hash string `gorm:"size:432;column:hash"`
}

// TableName Set User's table name to be `users`
func (User) TableName() string {
	return "users"
}

// UserGetList gets a map of all `User`s
func UserGetList() map[int64]User {
	var users []User
	usermap := make(map[int64]User)
	db, err := iniC("error(user-get-list) loading database\n")
	defer db.Close()
	if !err {
		db.Find(&users)
		fmt.Printf("- found %d entries\n", len(users))
		for _, x := range users {
			usermap[x.ID] = x
		}
	}
	return usermap
}

/* http://jinzhu.me/gorm/crud.html#query */

// ByName gets a user by [name].
// If `u` properties are set, then those are defaulted in FirstOrInit
func (u *User) ByName(name string) bool {
	db, err := iniC("error(user-by-name) loading database\n")
	defer db.Close()
	if !err {
		db.FirstOrInit(u, User{Name: name})
	}
	return u.Name == name
}

// ByID gets a user by [id].
func (u *User) ByID(id int64) bool {
	db, err := iniC("error(user-by-id) loading database\n")
	defer db.Close()
	if !err {
		db.FirstOrInit(u, User{ID: id})
	}
	return u.ID == id
}

// CreateSession32 creates sessioni with a default salt size.
func (u *User) CreateSession32(r interface{}, hours int, host string) (bool, Session) {
	return u.CreateSession(r, hours, host, 32)
}

// CreateSession is a test to attempt to save a session into the sessions table.
//
// FIXME: we should be checking if there is a existing record in sessions table
// and re-using it for the user executing UPDATE as opposed to CREATE.
func (u *User) CreateSession(r interface{}, hours int, host string, saltSize int) (bool, Session) {

	t := time.Now()
	salt := util.NewSaltString(saltSize)
	result := false

	sess := Session{
		Host:    host,
		UserID:  u.ID,
		SessID:  salt,
		Created: t,
		Expires: t.Add(durationHrs(hours)),
	}

	switch d := r.(type) {
	case *http.Response:
		sess.Client = util.ToBase64(d.Request.RemoteAddr)
	case string:
		sess.Client = util.ToBase64(d)
	default:
		sess.Client = util.ToBase64(unknownclient)
	}

	db, err := u.iniC("error(create-session) loading database\n")

	defer db.Close()
	if !err {
		db.Create(&sess)
		if db.RowsAffected == 1 {
			result = true
		}
	}

	return result, sess
}

// Create attempts to create a user and returns success or failure.
// If a user allready exists results in failure.
//
// Returns
// (-1) `db.Open`,
// (1) `User.Name` exists
// (0) on success
func (u *User) Create(name string, pass string, saltSize int) int {

	mysalt := saltsize
	if saltSize != -1 {
		mysalt = saltSize
	}

	db, err := u.iniC("error(user-create): loading database\n")
	if err {
		return -1
	}

	tempUser := User{ID: -1}
	db.FirstOrInit(&tempUser, User{Name: name})
	if tempUser.ID != -1 {
		db.Close()
		return 1 // user exists
	}

	bsalt := util.NewSaltCSRNG(mysalt)
	u.Name = name
	u.Salt = util.BytesToBase64(bsalt)
	u.Hash = util.BytesToBase64(util.GetPasswordHash(pass, bsalt))

	defer db.Close()
	db.Create(u)

	return 0
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

// ValidatePassword checks against a provided salt and hash.
// we use the user's [name] to find the table-record
// and then validate the password.
func (u *User) ValidatePassword(pass string) bool {
	// open the database
	db, err := u.iniC("error(validate-password) loading database\n")
	if err {
		return false
	}

	result := false
	tempUser := User{Name: "really?"}
	db.FirstOrInit(&tempUser, User{Name: u.Name})

	if db.RowsAffected == 0 && tempUser.Name != u.Name {
		db.Close()
		fmt.Println("Record not found")
		return false
	}

	defer db.Close()
	if tempUser.Name != u.Name {
		fmt.Printf("- no user found. %v\n", tempUser)
	} else {
		result = tempUser.validate(pass)
	}

	return result
}

// ValidateSession checks against a provided salt and hash.
// BUT FIRST, it checks for a valid session?
func (u *User) ValidateSession(host string, client interface{}) bool {

	db, err := u.iniC("error(validate-session) loading database\n")
	if err {
		return false
	}

	clistr := ""
	// cess := ""
	switch c := client.(type) {
	case *http.Response:
		clistr = util.ToUBase64(c.Request.RemoteAddr)
	case string:
		clistr = util.ToUBase64(c)
	default:
		clistr = util.ToBase64(unknownclient)
	}

	var sess Session
	db.Where("[client] = ? AND [host] = ?", clistr, host).FirstOrInit(&sess, Session{ID: -1})
	defer db.Close()
	if sess.ID == -1 {
		return false
	}

	t := time.Now()
	if t.After(sess.Expires) {
		return false
	}

	return true
}

// EnsureTableUsers creates table [users] if not exist.
func EnsureTableUsers() {
	var u User
	db, _ := iniK("error(ensure-table-users) loading db (perhaps expected)\n")
	// if !e {
	defer db.Close()
	if !db.HasTable(u) {
		db.CreateTable(u)
	}
	// }
}
