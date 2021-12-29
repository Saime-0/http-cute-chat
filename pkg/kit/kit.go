package kit

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func IntPtr(i int) *int              { return &i }
func StringPtr(s string) *string     { return &s }
func Ptr(i interface{}) *interface{} { return &i }

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

// deprecated
func SyntSQL(i interface{}) string {
	return match(i)
}
func match(v interface{}) string {
	switch v.(type) {
	case string:
		return "'" + v.(string) + "'"
	case bool:
		return strings.ToUpper(strconv.FormatBool(v.(bool)))
	case int:
		return strconv.Itoa(v.(int))
	case int64:
		return strconv.Itoa(int(v.(int64)))
	case *int64:
		if v.(*int64) != nil {
			return match(*v.(*int64))
		}
	case *bool:
		if v.(*bool) != nil {
			return match(*v.(*bool))
		}
	case *int:
		if v.(*int) != nil {
			return match(*v.(*int))
		}
	case *string:
		if v.(*string) != nil {
			return match(*v.(*string))
		}
	default:
		return fmt.Sprintf("'%v'", v)
	}
	return "NULL"
}

func IntSQLArray(arr []int) string {
	sqlArr := ""
	for _, v := range arr {
		//switch v.(type) {
		//case string, rune:
		//	v = spew.Sprint("'", v, "'")
		//}
		sqlArr += fmt.Sprint(",", v)
	}
	return "(" + TrimFirstRune(sqlArr) + ")"
}
