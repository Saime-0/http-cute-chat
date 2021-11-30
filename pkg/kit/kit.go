package kit

import "unicode/utf8"

func IntPtr(i int) *int { return &i }

func LeastOne(args ...bool) (discover bool) {
	for _, arg := range args {
		if arg {
			return true
		}
	}
	return
}

func TrimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
