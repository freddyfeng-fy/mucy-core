package numbers

import (
	"encoding/json"
	"strconv"
)

// ToInt64 str to int64，如果转换失败，默认值为0
// str 字符串
func ToInt64(str string) int64 {
	return ToInt64ByDefault(str, 0)
}

// ToInt64ByDefault str to int64
// str 字符串
// def 如果转换失败使用的默认值
func ToInt64ByDefault(str string, def int64) int64 {
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		val = def
	}
	return val
}

// ToInt str to int，如果转换失败，默认值为0
// str 字符串
func ToInt(str string) int {
	return ToIntByDefault(str, 0)
}

// ToIntByDefault str to int
// str 字符串
// def 如果转换失败使用的默认值
func ToIntByDefault(str string, def int) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		val = def
	}
	return val
}

func GetInterfaceToInt(t1 interface{}) int {
	var t2 int
	switch t1.(type) {
	case uint:
		t2 = int(t1.(uint))
		break
	case int8:
		t2 = int(t1.(int8))
		break
	case uint8:
		t2 = int(t1.(uint8))
		break
	case int16:
		t2 = int(t1.(int16))
		break
	case uint16:
		t2 = int(t1.(uint16))
		break
	case int32:
		t2 = int(t1.(int32))
		break
	case uint32:
		t2 = int(t1.(uint32))
		break
	case int64:
		t2 = int(t1.(int64))
		break
	case uint64:
		t2 = int(t1.(uint64))
		break
	case float32:
		t2 = int(t1.(float32))
		break
	case float64:
		t2 = int(t1.(float64))
		break
	case string:
		t2, _ = strconv.Atoi(t1.(string))
		if t2 == 0 && len(t1.(string)) > 0 {
			f, _ := strconv.ParseFloat(t1.(string), 64)
			t2 = int(f)
		}
		break
	case nil:
		t2 = 0
		break
	case json.Number:
		t3, _ := t1.(json.Number).Int64()
		t2 = int(t3)
		break
	default:
		t2 = t1.(int)
		break
	}
	return t2
}
