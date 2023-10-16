package strs

import (
	uuid "github.com/satori/go.uuid"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

/*
IsBlank checks if a string is whitespace or empty (""). Observe the following behavior:

	goutils.IsBlank("")        = true
	goutils.IsBlank(" ")       = true
	goutils.IsBlank("bob")     = false
	goutils.IsBlank("  bob  ") = false

Parameter:

	str - the string to check

Returns:

	true - if the string is whitespace or empty ("")
*/
func IsBlank(str string) bool {
	strLen := len(str)
	if str == "" || strLen == 0 {
		return true
	}
	for i := 0; i < strLen; i++ {
		if unicode.IsSpace(rune(str[i])) == false {
			return false
		}
	}
	return true
}

func IsNotBlank(str string) bool {
	return !IsBlank(str)
}

func IsAnyBlank(strs ...string) bool {
	for _, str := range strs {
		if IsBlank(str) {
			return true
		}
	}
	return false
}

func DefaultIfBlank(str, def string) string {
	if IsBlank(str) {
		return def
	} else {
		return str
	}
}

// IsEmpty checks if a string is empty (""). Returns true if empty, and false otherwise.
func IsEmpty(str string) bool {
	return len(str) == 0
}

func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}

func Substr(s string, start, length int) string {
	bt := []rune(s)
	if start < 0 {
		start = 0
	}
	if start > len(bt) {
		start = start % len(bt)
	}
	var end int
	if (start + length) > (len(bt) - 1) {
		end = len(bt)
	} else {
		end = start + length
	}
	return string(bt[start:end])
}

func Equals(a, b string) bool {
	return a == b
}

func EqualsIgnoreCase(a, b string) bool {
	return a == b || strings.ToUpper(a) == strings.ToUpper(b)
}

// RuneLen 字符成长度
func RuneLen(s string) int {
	bt := []rune(s)
	return len(bt)
}

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Int64ToStr(num int64) string {
	return strconv.FormatInt(num, 10)
}

func IntToStr(num int) string {
	return strconv.Itoa(num)
}

func UUID() string {
	u := uuid.NewV4()
	return strings.ReplaceAll(u.String(), "-", "")
}
