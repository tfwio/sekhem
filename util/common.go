package util

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const (
	resErrorCacheFile = "- error: cache file [%s]\n"
)

// FileExists checks if a given path or file exists.
// I don't believe this is aware if this is a file or not.
func FileExists(pFilePath string) bool {
	if _, err := os.Stat(pFilePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// Touch will create a file if it does not exist returns success.
// Bare in mind if the file exists before calling, then this will return false.
func Touch(pFilePath string) bool {

	if FileExists(pFilePath) {
		return false
	}
	var file, err = os.Create(pFilePath)
	defer file.Close()
	if err != nil {
		return false
	}
	return true
}

// GetDirectory expects a file as input and returns
// its parent directory.
// if input is a directory, I'm wondering what happens.
func GetDirectory(pFilePath string) (string, error) {
	dir, err := filepath.Abs(filepath.Dir(pFilePath))
	return dir, err
}

// StripFileExtension ...yep.
func StripFileExtension(pFilePath string) string {
	return strings.Replace(pFilePath, filepath.Ext(pFilePath), "", -1)
}

// Abs returns an absolute representation of path; Ignores errors.
func Abs(pPath string) (dir string) {
	dir, _ = filepath.Abs(pPath)
	return dir
}

// CacheFile Loads a local file in to string
func CacheFile(pFilePath string) string {
	mop, err := ioutil.ReadFile(pFilePath)
	if err != nil {
		return fmt.Sprintf(resErrorCacheFile, pFilePath)
	}
	return string(mop)
}

// StrInt64 string to int helper
func StrInt64(pStrInput string) int64 {
	var err error
	var fpoop int64
	if fpoop, err = strconv.ParseInt(pStrInput, 10, 32); err != nil {
		return 0
	}
	return int64(fpoop)
}

// TrimUnixSlash trims left and right forward-slashes from input string.
func TrimUnixSlash(pStrInput string) string {
	return TrimUnixSlashRight(TrimUnixSlashLeft(pStrInput))
}

// TrimUnixSlashLeft trims leftmost forward-slash from input string.
func TrimUnixSlashLeft(pStrInput string) string {
	return strings.TrimLeft(pStrInput, "/")
}

// TrimUnixSlashRight trims right forward-slash from input string.
func TrimUnixSlashRight(pStrInput string) string {
	return strings.TrimRight(pStrInput, "/")
}

// Space2Dash converts or replaces all spaces in a string with a small-dash.
func Space2Dash(pStrInput string) string {
	return strings.Replace(pStrInput, " ", "-", -1)
}

// UnixSlash converts all backslash to forward-slash.
func UnixSlash(instr string) string {
	return strings.Replace(instr, "\\", "/", -1)
}

// OSSlash converts all backslash to forward-slash (if OS is not windows).
// It'd probably be best to just use your standard `fileutil.Abs(…)`.
func OSSlash(instr string) string {
	if runtime.GOOS == "windows" {
		return strings.Replace(instr, "\\", "/", -1)
	}
	return instr
}

// StrTransformLiteral takes a literal string stuff like `EOL` and asserts
// literal code(s) such as `\n` and takes a measure or two to clean up the
// string to something a bit more normative.  You might say, this makes
// a string suitable for a JSON value ---might not.
func StrTransformLiteral(input string) (str string) {
	str = strings.Replace(input, `\r\n`, "\n", -1)
	str = strings.Replace(str, `\n`, "\n", -1)
	str = strings.Replace(str, `\t`, "	", -1)
	str = strings.Replace(str, `\\`, `\`, -1)
	str = strings.Replace(str, `\"`, `"`, -1)
	str = strings.Trim(str, "\"")
	return
}

// ConvertTransient What does this actually do?
func ConvertTransient(pInput string) string {
	return Space2Dash(
		TrimUnixSlash(
			UnixSlash(
				filepath.Dir(pInput))))
}

//func shaFile() {
//	s := "sha1 me"
//	h := sha1.New()
//	h.Write([]byte(s))
//	bs := h.Sum(nil)
//	fmt.Println(s)
//	fmt.Printf("%x\n", bs)
//}

// Sha1String just gets SHA1.
func Sha1String(pStrData string) string {
	hasher := sha1.New()
	hasher.Write([]byte(pStrData))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

// UNUSED
// func sha1Bytes(pStrData string) []byte {
// 	hasher := sha1.New()
// 	hasher.Write([]byte(pStrData))
// 	return hasher.Sum(nil)
// }

// CatArrayPad - Concatenate a string
// were padding the buffer here with a single char.
func CatArrayPad(pStrArray []string, pad string) string {
	if len(pStrArray) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for _, str := range pStrArray {
		buffer.WriteString(str + pad)
	}
	return strings.Trim(buffer.String(), pad) // fmt.Println(buffer.String())
}

// CatArray - Concatenate a string, or empty string.
func CatArray(input []string) string {
	if len(input) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for _, str := range input {
		buffer.WriteString(str)
	}
	return buffer.String() // fmt.Println(buffer.String())
}

// Cat - Concatenate a string by way of writing input to a buffer and
// converting returning its .WriteString() function.
func Cat(pInputString ...string) string {
	var buffer bytes.Buffer
	for _, str := range pInputString {
		buffer.WriteString(str)
	}
	return buffer.String() // fmt.Println(buffer.String())
}

// Insert inserts the value into the slice at the specified index,
// which must be in range.
// The slice must have room for the new element.
func Insert(slice []int, index, value int) []int {
	// Grow the slice by one element.
	slice = slice[0 : len(slice)+1]
	// Use copy to move the upper part of the slice out of the way and open a hole.
	copy(slice[index+1:], slice[index:])
	// Store the new value.
	slice[index] = value
	// Return the result.
	return slice
}

// CharIsNumber checks wether input string contains all digit characters.
func CharIsNumber(input string) bool {
	for _, b := range []byte(input) {
		if !(b >= 48 && b <= 57) {
			return false
		}
	}
	return true
}

const unknownString = "unknown date"

// CheckDateString checks the beginning of a file-name for an 8-digit date-string;
// I.E.: `YYYYMMDD`
func CheckDateString(input string) string {
	result := strings.Index(input, " ")
	// println(fmt.Sprintf("first-index:  %d", result))
	if result >= 0 && result == 8 && CharIsNumber(input[:8]) {
		return input[:8]
	}
	return unknownString
}
