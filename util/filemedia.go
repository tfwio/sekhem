package util

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
)

const mediaExtensions = ".aac,.alac,.dsf,.flac,.m4a,.m4b,.m4p,.mp4,.mp3,.ogg"
const mediaTagExtensions = ".aac,.alac,.dsf,.flac,.m4a,.m4b,.m4p,.mp4,.mp3,.ogg"

func IsMediaFile(path string) bool {
	ref := strings.ToLower(filepath.Ext(path))
	for _, ext := range strings.Split(mediaExtensions, ",") {
		if strings.ToLower(ext) == ref {
			return true
		}
	}
	return false
}
func GetMediaFile(path string) (tag.Metadata, error) {
	if !IsMediaFile(path) {
		return nil, errors.New("Not media")
	}
	file, ferr := os.Open(path)
	if ferr == nil {
		meta, terr := tag.ReadFrom(file)
		file.Close()
		return meta, terr
	}
	return nil, ferr

}
