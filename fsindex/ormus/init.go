package ormus

import (
	"fmt"

	"github.com/jinzhu/gorm"
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
