// +build !windows

package util

// https://stackoverflow.com/questions/20875336/how-can-i-get-a-files-ctime-atime-mtime-and-change-them-using-golang

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func statTimes(name string) (atime, mtime, ctime time.Time, err error) {
	fi, err := os.Stat(name)
	if err != nil {
		return
	}
	mtime = fi.ModTime()
	stat := fi.Sys().(*syscall.Stat_t)
	atime = time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))
	ctime = time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))
	return
}

// !
func nonono() {
	name := "stat.file"
	atime, mtime, ctime, err := statTimes(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(atime, mtime)
	fmt.Println(ctime)
	err = os.Chtimes(name, atime, mtime)
	if err != nil {
		fmt.Println(err)
		return
	}
	atime, mtime, ctime, err = statTimes(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(atime, mtime)
	fmt.Println(ctime)
}
