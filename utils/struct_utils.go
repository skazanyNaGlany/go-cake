package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/ssrathi/go-attr"
)

type StructUtils struct{}
type TagMap map[string]map[string]string

var StructUtilsInstance StructUtils

func (su StructUtils) FormatStruct(val any) (string, error) {
	formatted := ""

	values, err := attr.Values(val)

	if err != nil {
		return "", err
	}

	formatted += fmt.Sprintf("[%T]\n", val)

	for name, value := range values {
		formatted += fmt.Sprintf("%v: %v\n", name, value)
	}

	return strings.TrimSpace(formatted), nil
}

func (su StructUtils) FormatStructNoError(val any) string {
	formatted, _ := su.FormatStruct(val)

	return formatted
}

func (su StructUtils) GetCleanType(val any) string {
	typeStr := fmt.Sprintf("%T", val)

	parts := strings.Split(typeStr, ".")
	len_parts := len(parts)

	if len_parts > 1 {
		return parts[len_parts-1]
	}

	return typeStr
}

func (su StructUtils) StructToJSONString(s any) (string, error) {
	str, err := json.Marshal(s)

	if err != nil {
		return "", err
	}

	return string(str), nil
}

func (su StructUtils) JSONStringToMap(s string) (map[string]any, error) {
	var result map[string]any

	decoder := json.NewDecoder(bytes.NewReader([]byte(s)))
	decoder.UseNumber()

	err := decoder.Decode(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (su StructUtils) JSONStringToArray(s string) ([]any, error) {
	result := make([]any, 0)

	if err := json.Unmarshal([]byte(s), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (su StructUtils) StructToMap(s any) (map[string]any, error) {
	itemBytes, err := json.Marshal(s)

	if err != nil {
		return nil, err
	}

	jsonObjectMap := make(map[string]any)

	err = json.Unmarshal(itemBytes, &jsonObjectMap)

	return jsonObjectMap, err
}

func (su StructUtils) StructToTagMap(
	s any,
	neededTags []string,
	key string) (TagMap, error) {
	tags := make(TagMap)

	fieldNames, err := attr.Names(s)

	if err != nil {
		return nil, err
	}

	for _, ifieldName := range fieldNames {
		itag := make(map[string]string)

		for _, iNeededTag := range neededTags {
			if iNeededTag == "name" {
				itag["name"] = ifieldName
				continue
			}

			tagValue, err := attr.GetTag(s, ifieldName, iNeededTag)

			if err != nil {
				return tags, err
			}

			tagValue = strings.TrimSpace(tagValue)
			tagValueParts := strings.Split(tagValue, ",")

			if len(tagValueParts) > 1 {
				tagValue = strings.TrimSpace(tagValueParts[0])
			}

			itag[iNeededTag] = tagValue
		}

		if key == "name" {
			tags[ifieldName] = itag
		} else {
			if itag[key] != "" {
				tags[itag[key]] = itag
			}
		}
	}

	return tags, nil
}

func (su *StructUtils) GetFinalValue(s any) any {
	kind := reflect.TypeOf(s).Kind()

	if kind != reflect.Pointer {
		return s
	}

	valueOf := reflect.ValueOf(s)

	if valueOf.Pointer() == 0 {
		return s
	}

	return valueOf.Elem().Interface()
}
