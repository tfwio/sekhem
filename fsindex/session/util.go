package session

import "encoding/base64"

// FromBase64e gets base-64 StdEncoding (with error)
func fromBase64e(input string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(input)
}

// FromBase64 gets base-64 StdEncoding; ignores error.
func fromBase64(input string) []byte {
	result, _ := base64.StdEncoding.DecodeString(input)
	return result
}

// ToBase64 gets base-64 StdEncoding
func toBase64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

// ToUBase64 gets base-64 URLEncoding
func toUBase64(input string) string {
	return base64.URLEncoding.EncodeToString([]byte(input))
}

// FromUBase64 gets base-64 URLEncoding
func fromUBase64(input string) string {
	result, _ := base64.URLEncoding.DecodeString(input)
	return string(result)
}

// BytesToBase64 gets base-64 StdEncoding
func bytesToBase64(input []byte) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}
