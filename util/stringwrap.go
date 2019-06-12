package util

import "strings"

// Wrap wraps text with `wrap`, written for converting "v" to "/v/".
// see: https://blog.golang.org/strings
func Wrap(text string, wrap string) string {
	result := text
	if strings.Index(result, wrap) != 0 {
		result = Cat(wrap, result)
	}
	if strings.LastIndex(result, wrap) != (len(result) - 1) {
		result = Cat(result, wrap)
	}
	return result
}

// Wrapper concatenates text and wraps it like `Wrap` does with `sep`-arator.
func Wrapper(sep string, text ...string) string {
	return Wrap(strings.Join(text, sep), sep)
}

// WrapperRight concatenates text and wraps it like `Wrap` and pads it to the right with `sep`.
func WrapperRight(sep string, text ...string) string {
	return WrapRight(strings.Join(text, sep), sep)
}

// WrapperLeft concatenates text and wraps it like `Wrap` and pads it to the left with `sep`.
func WrapperLeft(sep string, text ...string) string {
	return WrapperLeft(strings.Join(text, sep), sep)
}

// WReap ensures text that should be wrapped is re-wrapped.
func WReap(text string, wrap string) string {
	return Wrap(strings.Trim(text, wrap), "/")
}

// WReaper ensures text that should be wrapped is re-wrapped.
func WReaper(wrap string, text ...string) string {
	for i, t := range text {
		text[i] = strings.Trim(t, wrap)
	}
	return Wrap(Wrapper(wrap, text...), wrap)
}

// WReaperRight ensures text that should be wrapped is re-wrapped.
func WReaperRight(wrap string, text ...string) string {
	for i, t := range text {
		text[i] = strings.Trim(t, wrap)
	}
	return WrapRight(Wrapper(wrap, text...), wrap)
}

// WReaperLeft ensures text that should be wrapped is re-wrapped.
func WReaperLeft(wrap string, text ...string) string {
	for i, t := range text {
		text[i] = strings.Trim(t, wrap)
	}
	return WrapLeft(Wrapper(wrap, text...), wrap)
}

// WrapLeft puts `wrap` at the beginning of the string if not already present.
func WrapLeft(wrap string, text string) string {
	result := text
	if strings.Index(result, wrap) != 0 {
		result = Cat(wrap, result)
	}
	return result
}

// WrapRight puts `wrap` at the end of the string if not already present.
func WrapRight(wrap string, text string) string {
	result := text
	if strings.LastIndex(result, wrap) != (len(result) - 1) {
		result = Cat(result, wrap)
	}
	return result
}
