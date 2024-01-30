package utils

import (
	"fmt"
	"reflect"

	"github.com/thoas/go-funk"
)

type IterableUtils struct{}

var IterableUtilsInstance IterableUtils

func (iu IterableUtils) IterableHasUniqueItems(iterable any) (bool, error) {
	iterableType := reflect.TypeOf(iterable)

	if (ReflectUtils{}).IsArrayType(iterableType) {
		items := make([]any, 0)
		unique := true

		funk.ForEach(iterable, func(iitem any) {
			if !unique {
				return
			}

			if funk.Contains(items, iitem) {
				unique = false
			} else {
				items = append(items, iitem)
			}
		})

		return unique, nil
	} else if (ReflectUtils{}).IsMapType(iterableType) {
		keys := make([]any, 0)
		unique := true

		funk.ForEach(iterable, func(iikey, ivalue any) {
			if !unique {
				return
			}

			if funk.Contains(keys, iikey) {
				unique = false
			} else {
				keys = append(keys, iikey)
			}
		})

		return unique, nil
	}

	return false, fmt.Errorf("'%T' is not iterable", iterable)
}
