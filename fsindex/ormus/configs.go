package ormus

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

var (
	datasource              string
	datasys                 string
	saltsize                = 48
	defaultSessionLength, _ = time.ParseDuration("2h")
	unknownclient           = "unknown-client"
)

// this is just an empty struct to be inherited for simple loading
// of orm/db.
type dbx struct {
}

// close on error (no message)
func dbinit() (*gorm.DB, bool) {
	db, err := gorm.Open(datasys, datasource)
	result := false
	if err != nil {
		db.Close()
		result = true
	}
	return db, result
}

// no close on error (no message)
func dbinik() (*gorm.DB, bool) {
	db, err := gorm.Open(datasys, datasource)
	result := false
	if err != nil {
		result = true
	}
	return db, result
}

// close on error
func (dbx) init() (*gorm.DB, bool) {
	return dbinit()
}

// close on error
func (dbx) iniC(format string, msg ...interface{}) (*gorm.DB, bool) {
	return inik(true, format, msg...)
}

// keep on error
func (dbx) iniK(format string, msg ...interface{}) (*gorm.DB, bool) {
	return inik(false, format, msg...)
}

// closes the database and prints requested status on error.
func iniC(format string, msg ...interface{}) (*gorm.DB, bool) {
	return inik(true, format, msg...)
}

// keep error
func iniK(format string, msg ...interface{}) (*gorm.DB, bool) {
	return inik(false, format, msg...)
}

// closes the database and prints requested status on error.
func inik(closeOnError bool, format string, msg ...interface{}) (*gorm.DB, bool) {
	db, e := dbinik()
	if e {
		if format != "" {
			fmt.Printf(format, msg...)
		}
		if closeOnError {
			db.Close()
		}
	}
	return db, e
}

// SetDefaults allows a external library to set the local datasource.
// Set saltSize to -1 to persist default.
func SetDefaults(source string, sys string, saltSize int) {
	datasource = source
	datasys = sys
	if saltSize != -1 {
		saltsize = saltSize
	}
}

// returns calculated duration or on error the default session length '2hr'
func durationHrs(hr int) time.Duration {
	if result, err := time.ParseDuration(fmt.Sprintf("%vh", hr)); err != nil {
		return result
	}
	return defaultSessionLength
}

// FIXME: This IS NOT USED
func loginCheckUP(user string, pass string) bool {
	result := false
	if len(user) > 4 && len(pass) > 3 {
		u := User{Name: user}
		result = u.ValidatePassword(pass)
		//fmt.Printf("Result: %v \n", result)
	} else {
		//fmt.Printf("- username %s; pass %s\n", user, pass)
		println("- username must be > len(3) chars long")
		println("- you must supply a password > len(4) chars long")
	}
	return result
}
