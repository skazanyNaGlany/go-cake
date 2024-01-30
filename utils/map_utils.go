package utils

import (
	"reflect"

	"github.com/thoas/go-funk"
)

type MapUtils struct{}

var MapUtilsInstance MapUtils

func (mu MapUtils) GetMapStringKeys(m map[string]any, recursive bool) []string {
	keysLocal := funk.Keys(m).([]string)

	if recursive {
		for _, v := range m {
			if reflect.ValueOf(v).Kind() != reflect.Map {
				continue
			}

			vCasted, ok := v.(map[string]any)

			if !ok {
				continue
			}

			keysLocal = append(keysLocal, mu.GetMapStringKeys(vCasted, recursive)...)
		}
	}

	return keysLocal
}

func (mu MapUtils) MapStringCopy(m map[string]any, skipFields []string) map[string]any {
	copy := make(map[string]any, 0)

	for k, v := range m {
		if funk.ContainsString(skipFields, k) {
			continue
		}

		copy[k] = v
	}

	return copy
}
