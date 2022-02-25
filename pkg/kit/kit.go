package kit

import (
	"encoding/base32"
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"
	"unicode/utf8"
)

func IntPtr(i int) *int              { return &i }
func Int64Ptr(i int64) *int64        { return &i }
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

func IntSQLArray(arr []int) string {
	sqlArr := ""
	for _, v := range arr {
		sqlArr += fmt.Sprint(",", v)
	}
	return "(" + TrimFirstRune(sqlArr) + ")"
}

func RandomSecret(length int) string {
	rand.Seed(time.Now().UnixNano())
	letterRunes := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567")

	bytes := make([]rune, length)

	for i := range bytes {
		bytes[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(bytes)
}

func CringeSecret(blength int) string {
	var result string
	secret := make([]byte, blength)
	gen, err := rand.Read(secret)
	if err != nil || gen != blength {
		return result
	}
	var encoder = base32.StdEncoding.WithPadding(base32.NoPadding)
	result = encoder.EncodeToString(secret)
	return result
}

func GetUniqueInts(arr []int) []int {
	keys := make(map[int]bool)
	var list []int
	for _, entry := range arr {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func After(offset int64) int64 {
	return time.Now().Unix() + offset
}

func Err(v interface{}, err error) error {
	return err
}

var IsLetter = regexp.MustCompile(`^[a-zA-Z0-9\-\=]{20}$`).MatchString
