package arrays

import (
	"reflect"
	"strings"
)

func Contains(search interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == search {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(search)).IsValid() {
			return true
		}
	}
	return false
}

func ContainsIgnoreCase(search string, target []string) bool {
	if len(search) == 0 {
		return false
	}
	if len(target) == 0 {
		return false
	}
	search = strings.ToLower(search)
	for i := 0; i < len(target); i++ {
		if strings.ToLower(target[i]) == search {
			return true
		}
	}
	return false
}

func SplitSlice(input interface{}, segLength int) interface{} {
	if segLength < 1 {
		segLength = 1 //传入0会导致死循环
	}

	sliceOfElemSlice := reflect.SliceOf(reflect.TypeOf(input))
	listOfElem := reflect.New(sliceOfElemSlice)

	varValue := reflect.ValueOf(input)
	inputLen := varValue.Len()
	for start := 0; start < inputLen; {
		end := start + segLength
		if end >= inputLen {
			end = inputLen
		}
		listOfElem.Elem().Set(reflect.Append(listOfElem.Elem(), varValue.Slice(start, end)))
		start = end
	}
	return listOfElem.Interface()
}

func containsDuplicate(nums []int) bool {
	n := make(map[int]int)
	for i := 0; i < len(nums); i++ {
		if _, ok := n[nums[i]]; ok {
			return true
		}
		n[nums[i]] = i
	}
	return false
}
