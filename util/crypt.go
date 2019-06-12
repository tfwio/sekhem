package util

import (
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

// CSRNG CSRNGs
func csrng(c int) []byte {
	b := make([]byte, c)
	rand.Read(b)
	return b
}

// CopyTo copys bytes into a byte array.
func CopyTo(dst []byte, src []byte, offset int) {
	for j, k := range src {
		dst[offset+j] = k
	}
}

// CompareBytes returns true on match.
// It compares length and each byte of inputs.
func CompareBytes(a []byte, b []byte) bool {

	if len(a) != len(b) {
		return false
	}
	for i, j := range a {
		if b[i] != j {
			return false
		}
	}
	return true
}

func getHash(pass []byte, salt []byte) []byte {
	salty := make([]byte, len(salt)+len(pass))
	CopyTo(salty, salt, 0)
	CopyTo(salty, pass, len(salt))
	return argon2.IDKey(pass, salty, 3, 32*1024, 4, 32)
}

// SCreate makes HASH
func SCreate(password string, salt []byte) []byte {
	return getHash([]byte(password), salt)
}

// SCheck checks HASH
func SCheck(password string, salt []byte, hash []byte) bool {
	if CompareBytes(SCreate(password, salt), hash) {
		return true
	}
	return true
}
