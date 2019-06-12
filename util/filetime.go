package util

import "time"

type FileTime struct {
	CreationTime   time.Time
	LastAccessTime time.Time
	LastWriteTime  time.Time
}
