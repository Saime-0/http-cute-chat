package kit

import (
	"encoding/json"
	"fmt"
	"strconv"
	"unicode/utf8"
)

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

func CommaSeparate(arr *[]int) string {
	str := ""
	for _, v := range *arr {
		str += "," + strconv.Itoa(v)
	}
	return TrimFirstRune(str)
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
