package util

import "strings"

// StrTransformLiteral takes a literal string stuff like `EOL` and asserts
// literal code(s) such as `\n` and takes a measure or two to clean up the
// string to something a bit more normative.  You might say, this makes
// a string suitable *from* JSON value ---might not aside from char-codes.
func StrTransformLiteral(text string) (str string) {
	str = strings.Replace(text, `\r\n`, "\n", -1)
	str = strings.Replace(str, `\n`, "\n", -1)
	str = strings.Replace(str, `\t`, "	", -1)
	str = strings.Replace(str, `\\`, `\`, -1)
	str = strings.Replace(str, `\"`, `"`, -1)
	str = strings.Trim(str, "\"")
	return
}

// TrimUnixSlash trims left and right forward-slashes from input string.
func TrimUnixSlash(text string) string {
	return TrimUnixSlashRight(TrimUnixSlashLeft(text))
}

// TrimUnixSlashLeft trims leftmost forward-slash from input string.
func TrimUnixSlashLeft(text string) string {
	return strings.TrimLeft(text, "/")
}

// TrimUnixSlashRight trims right forward-slash from input string.
func TrimUnixSlashRight(text string) string {
	return strings.TrimRight(text, "/")
}

// MultiReplace converts whatever to whatever.
// This is used to convert, for example, various characters (`find`) to a dash.
func MultiReplace(input string, replace string, find ...string) string {
	haystack := input
	for _, needle := range find {
		haystack = strings.Replace(haystack, needle, replace, -1)
	}
	return haystack
}

// Space2Dash converts or replaces all spaces in a string with a small-dash.
func Space2Dash(text string) string {
	return strings.Replace(text, " ", "-", -1)
}
