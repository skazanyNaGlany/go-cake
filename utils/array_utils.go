package utils

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"

	"github.com/thoas/go-funk"
	"github.com/yourbasic/radix"
)

type ArrayUtils struct{}

var ArrayUtilsInstance ArrayUtils

func (au ArrayUtils) ArrayPopString(slice []string, index int) ([]string, string) {
	item := slice[index]
	slice = append(slice[:index], slice[index+1:]...)

	return slice, item
}

func (au ArrayUtils) ArrayToMapArrayAny(slice []string) map[string][]any {
	result := make(map[string][]any)

	for _, item := range slice {
		result[item] = make([]any, 0)
	}

	return result
}

func (au ArrayUtils) Sort(s []any, ascending bool) ([]any, bool) {
	if len(s) <= 1 {
		return s, false
	}

	_type := reflect.TypeOf(s[0])

	if ReflectUtilsInstance.IsIntType(_type) {
		// array of signed integers
		sort.Slice(s, func(i, j int) bool {
			v1, _ := strconv.ParseInt(fmt.Sprintf("%v", s[i]), 10, 64)
			v2, _ := strconv.ParseInt(fmt.Sprintf("%v", s[j]), 10, 64)

			return v1 < v2
		})
	} else if ReflectUtilsInstance.IsUIntType(_type) {
		// array of unsigned integers
		sort.Slice(s, func(i, j int) bool {
			v1, _ := strconv.ParseUint(fmt.Sprintf("%v", s[i]), 10, 64)
			v2, _ := strconv.ParseUint(fmt.Sprintf("%v", s[j]), 10, 64)

			return v1 < v2
		})
	} else if ReflectUtilsInstance.IsFloatType(_type) {
		// array of floats
		sort.Slice(s, func(i, j int) bool {
			v1, _ := strconv.ParseFloat(fmt.Sprintf("%v", s[i]), 64)
			v2, _ := strconv.ParseFloat(fmt.Sprintf("%v", s[j]), 64)

			return v1 < v2
		})
	} else if ReflectUtilsInstance.IsComplexType(_type) {
		// array of complex
		sort.Slice(s, func(i, j int) bool {
			v1, _ := strconv.ParseComplex(fmt.Sprintf("%v", s[i]), 128)
			v2, _ := strconv.ParseComplex(fmt.Sprintf("%v", s[j]), 128)

			return real(v1) < real(v2)
		})
	} else if ReflectUtilsInstance.IsStringType(_type) {
		converted := make([]string, 0)
		converted2 := make([]any, 0)

		// convert to []string
		for _, value := range s {
			converted = append(converted, fmt.Sprintf("%v", value))
		}

		radix.Sort(converted)

		// convert back to []any
		for _, value2 := range converted {
			converted2 = append(converted2, value2)
		}

		s = converted2
	} else {
		return s, false
	}

	if !ascending {
		s = funk.Reverse(s).([]any)
	}

	return s, true
}
