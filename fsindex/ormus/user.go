package ormus

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tfwio/sekhem/util"
)

// User structure
type User struct {
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
	fmt.Printf("--> looking for %s\n", name)

	db, err := iniC("error(user-by-name) loading database\n")
	result := false
	defer db.Close()
	if !err {
		db.Where("[user] = ?", name).First(u)
		if u.Name == name {
			result = true
		}
	}
	fmt.Printf("!-> FOUND %s, %d\n", u.Name, u.ID)
	return result
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
func (u *User) CreateSession32(r *gin.Context, hours int, host string) (bool, Session) {
	return u.CreateSession(r, hours, host, 32)
}

// CreateSession is a test to attempt to save a session into the sessions table.
//
// FIXME: we should be checking if there is a existing record in sessions table
// and re-using it for the user executing UPDATE as opposed to CREATE.
func (u *User) CreateSession(r *gin.Context, hours int, host string, saltSize int) (bool, Session) {

	t := time.Now()
	ss := saltsize
	if saltSize != -1 {
		ss = saltSize
	}
	result := false
	sess := Session{
		Host:    host,
		UserID:  u.ID,
		SessID:  util.ToUBase64(util.NewSaltString(ss)),
		Created: t,
		Expires: t.Add(durationHrs(hours)),
	}

	sess.Client = getClientString(r)

	db, err := iniC("error(create-session) loading database\n")

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

	if u.ByName(name) {
		return 1
	}

	db, err := iniC("error(user-create): loading database\n")
	if err {
		return -1
	}

	bsalt := util.NewSaltCSRNG(mysalt)
	u.Name = name
	u.Salt = util.BytesToBase64(bsalt)
	u.Hash = util.BytesToBase64(util.GetPasswordHash(pass, bsalt))
	fmt.Printf("--> %s, %s, %v\n", u.Name, pass, u.Salt)

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
	db, err := iniC("error(validate-password) loading database\n")
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

// SessionRefresh Extend existing session.
// if succeeds, result is the created session, otherwise nil.
func (u *User) SessionRefresh(host string, client *gin.Context) interface{} {
	// println("--> SessionRefresh")
	// c := client.(*http.Request)
	clistr := getClientString(client)
	// fmt.Printf("cli-key: %s, host: %s\n", clistr, host)

	sess := Session{}
	db, err := iniC("error(validate-session) loading database\n")
	if err {
		return nil
	}
	db.LogMode(true)
	//
	db.First(&sess, "[cli-key] = ? AND [host] = ? AND [user_id] = ?", clistr, host, u.ID)
	//
	defer db.Close()
	if sess.UserID == u.ID {
		sess.Expires = time.Now().Add(durationHrs(2))
		db.Save(sess)
		// fmt.Printf("--> UPDATED SESS! user_id match %v, %v\n", sess.UserID, u.ID)
		return sess
	}
	// fmt.Printf("--> NO user_id match %v, %v\n", sess.UserID, u.ID)
	return nil
}

// UserSession checks a session for the user against the client/cookie.
// NOTE THAT USER ID MUST BE PRESENT!
func (u *User) UserSession(host string, client *gin.Context) (Session, bool) {
	clistr := getClientString(client)
	sess := Session{}
	db, err := iniC("error(validate-session) loading database\n")
	if err {
		return sess, false
	}
	db.LogMode(true)
	defer db.Close()
	db.First(&sess, "[cli-key] = ? AND [host] = ? AND [user_id] = ?", clistr, host, u.ID)
	return sess, db.RowsAffected != 0
}

// ValidateSessionByUserID checks against a provided salt and hash.
// BUT FIRST, it checks for a valid session?
func (u *User) ValidateSessionByUserID(host string, client *gin.Context) bool {

	println("--> ValidateSessionByUserID")

	clistr := getClientString(client)
	fmt.Printf("cli-key: %s, host: %s\n", clistr, host)

	result := false
	sess, err := u.UserSession(host, client)
	if err {
		fmt.Printf("--> Found no session\n")
		return false
	}

	if sess.UserID == u.ID {
		result = true
		fmt.Printf("--> user_id match %v, %v\n", sess.UserID, u.ID)
	} else {
		fmt.Printf("--> NO user_id match %v, %v\n", sess.UserID, u.ID)
	}

	t := time.Now()
	if !t.Before(sess.Expires) {
		fmt.Println("--> session isn't expired")
		result = false
	} else {
		fmt.Println("--> session is expired")
	}

	return result
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
