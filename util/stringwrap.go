package util

import "strings"

// Wrap wraps text with `wrap`, written for converting "v" to "/v/".
// see: https://blog.golang.org/strings
func Wrap(wrap string, text string) string {
	result := text
	if strings.Index(result, wrap) != 0 {
		result = Cat(wrap, result)
	}
	if strings.LastIndex(result, wrap) != (len(result) - 1) {
		result = Cat(result, wrap)
	}
	return result
}

// WrapLeft puts `wrap` at the beginning of the string if not already present.
func WrapLeft(separator string, text string) string {
	result := text
	if strings.Index(result, separator) != 0 {
		result = Cat(separator, result)
	}
	return result
}

// WrapRight puts `wrap` at the end of the string if not already present.
func WrapRight(separator string, text string) string {
	result := text
	if strings.LastIndex(result, separator) != (len(result) - 1) {
		result = Cat(result, separator)
	}
	return result
}

// Wrapper concatenates text and wraps it like `Wrap` does with `sep`-arator.
func Wrapper(separator string, text ...string) string {
	return Wrap(separator, strings.Join(text, separator))
}

// WrapperRight concatenates text and wraps it like `Wrap` and pads it to the right with `sep`.
func WrapperRight(separator string, text ...string) string {
	return WrapRight(strings.Join(text, separator), separator)
}

// WrapperLeft concatenates text and wraps it like `Wrap` and pads it to the left with `sep`.
func WrapperLeft(separator string, text ...string) string {
	return WrapperLeft(strings.Join(text, separator), separator)
}

// Trim trims each text element in the input array.
// TODO: should we act on `text` itself or return a new value?
func Trim(separator string, text ...string) []string {
	result := text
	for i, t := range result {
		result[i] = strings.Trim(t, separator)
	}
	return result
}

// TrimJoin trims each text element in the input array
// and `Join`s the result (using separator).
// TODO: should we act on `text` itself or return a new value?
func TrimJoin(separator string, text ...string) string {
	result := text
	for i, t := range result {
		result[i] = strings.Trim(t, separator)
	}
	return strings.Join(result, separator)
}

// WReap makes sure that each text node is trimmed of `separator` and also wraps text with the `separator`.
func WReap(separator string, text ...string) string {
	for i, t := range text {
		text[i] = strings.Trim(t, separator)
	}
	return Wrap(separator, Wrapper(separator, text...))
}

// WReapRight makes sure that each text node is trimmed of `separator` and also right-wraps text with the `separator`.
func WReapRight(separator string, text ...string) string {
	for i, t := range text {
		text[i] = strings.Trim(t, separator)
	}
	return WrapRight(Wrapper(separator, text...), separator)
}

// WReapLeft makes sure that each text node is trimmed of `separator` and also left-wraps text with the `separator`.
func WReapLeft(separator string, text ...string) string {
	for i, t := range text {
		text[i] = strings.Trim(t, separator)
	}
	return WrapLeft(Wrapper(separator, text...), separator)
}
