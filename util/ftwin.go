// +build windows

package util

// https://github.com/skillian/getfiletime/blob/master/getfiletime.go
// license: Unlicens

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

var (
	kernel32Handle          syscall.Handle
	fileTimeToSystemTimePtr uintptr
	getFileTime             uintptr
)

func init() {
	var err error
	kernel32Handle, err = syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		panic(fmt.Errorf("failed to load kernel32.dll: %v", err))
	}

	fileTimeToSystemTimePtr, err = syscall.GetProcAddress(
		kernel32Handle,
		"FileTimeToSystemTime")

	if err != nil {
		panic(fmt.Errorf(
			"failed to get FileTimeToSystemTime proc: %v", err))
	}

	getFileTime, err = syscall.GetProcAddress(
		kernel32Handle,
		"GetFileTime")

	if err != nil {
		panic(fmt.Errorf(
			"failed to get GetFileTime function from kernel32: %v",
			err))
	}
}

func createFileGenericRead(filename string) (h syscall.Handle, err error) {
	return syscall.CreateFile(
		syscall.StringToUTF16Ptr(filename),
		syscall.GENERIC_READ,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_NORMAL,
		0)
}

func fileTimeToSystemTime(ft *syscall.Filetime) (st syscall.Systemtime, err error) {
	res, _, err := syscall.Syscall(
		uintptr(fileTimeToSystemTimePtr),
		2,
		uintptr(unsafe.Pointer(ft)),
		uintptr(unsafe.Pointer(&st)),
		0)

	if err != syscall.Errno(0) {
		return syscall.Systemtime{}, fmt.Errorf(
			"failed to call FileTimeToSystemTime: %v", err)
	}

	if res == 0 {
		return syscall.Systemtime{}, fmt.Errorf(
			"call to FileTimeToSystemTime returned 0")
	}

	return st, nil
}

func GetFileTime(filename string) (FileTime, error) {
	fileHandle, err := createFileGenericRead(filename)
	if err != nil {
		return FileTime{}, fmt.Errorf(
			"failed to open file in GetFileTime: %v", err)
	}

	var fileCreate syscall.Filetime
	var fileAccess syscall.Filetime
	var fileModify syscall.Filetime

	res, _, err := syscall.Syscall6(
		uintptr(getFileTime),
		4,
		uintptr(fileHandle),
		uintptr(unsafe.Pointer(&fileCreate)),
		uintptr(unsafe.Pointer(&fileAccess)),
		uintptr(unsafe.Pointer(&fileModify)),
		0,
		0)

	if err != syscall.Errno(0) {
		return FileTime{}, fmt.Errorf(
			"failed to call GetFileTime: %v", err)
	}

	if res == 0 {
		return FileTime{}, fmt.Errorf(
			"Zero result from GetFileTime: %x", res)
	}

	if err := syscall.CloseHandle(fileHandle); err != nil {
		return FileTime{}, fmt.Errorf(
			"error closing file %v: %v",
			filename,
			err)
	}

	createSysTime, err := fileTimeToSystemTime(&fileCreate)
	if err != nil {
		return FileTime{}, fmt.Errorf(
			"failed to convert file create time (%v) to system time: %v",
			fileCreate,
			err)
	}
	accessSysTime, err := fileTimeToSystemTime(&fileAccess)
	if err != nil {
		return FileTime{}, fmt.Errorf(
			"failed to convert file access time (%v) to system time: %v",
			fileCreate,
			err)
	}
	modifySysTime, err := fileTimeToSystemTime(&fileModify)
	if err != nil {
		return FileTime{}, fmt.Errorf(
			"failed to convert file modify time (%v) to system time: %v",
			fileCreate,
			err)
	}

	return FileTime{
		CreationTime: time.Date(
			int(createSysTime.Year),
			time.Month(createSysTime.Month),
			int(createSysTime.Day),
			int(createSysTime.Hour),
			int(createSysTime.Minute),
			int(createSysTime.Second),
			int(createSysTime.Milliseconds)*int(time.Millisecond),
			time.UTC),
		LastAccessTime: time.Date(
			int(accessSysTime.Year),
			time.Month(accessSysTime.Month),
			int(accessSysTime.Day),
			int(accessSysTime.Hour),
			int(accessSysTime.Minute),
			int(accessSysTime.Second),
			int(accessSysTime.Milliseconds)*int(time.Millisecond),
			time.UTC),
		LastWriteTime: time.Date(
			int(modifySysTime.Year),
			time.Month(modifySysTime.Month),
			int(modifySysTime.Day),
			int(modifySysTime.Hour),
			int(modifySysTime.Minute),
			int(modifySysTime.Second),
			int(modifySysTime.Milliseconds)*int(time.Millisecond),
			time.UTC),
	}, nil
}
