package go_cake

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/skazanyNaGlany/go-cake/utils"
	attr "github.com/ssrathi/go-attr"
	"github.com/thoas/go-funk"
)

// DB model field specs
type FieldSpecs struct {
	Name            string
	JSON            string
	BSON            string
	Filterable      bool
	Projectable     bool
	Sortable        bool
	Insertable      bool
	Updatable       bool
	Hidden          bool
	Erased          bool
	ETag            bool
	RequireOnInsert bool
	RequireOnUpdate bool
	fieldType       reflect.Type
	availableSpecs  []string
}

func (fc FieldSpecs) TranslateField(field string) string {
	translatable := make(map[string]string)

	translatable["requireoninsert"] = "require-on-insert"
	translatable["requireonupdate"] = "require-on-update"

	if _, ok := translatable[field]; !ok {
		return field
	}

	return translatable[field]
}

func (fc FieldSpecs) GetKinds() (map[string]string, error) {
	filtered := make(map[string]string)

	kinds, err := attr.Kinds(fc)

	if err != nil {
		return filtered, err
	}

	for iname, ikind := range kinds {
		if iname == "Name" || iname == "JSON" || iname == "BSON" {
			continue
		}

		filtered[iname] = ikind
	}

	return filtered, err
}

func (fc *FieldSpecs) setSpecByKind(specName, value, kind string) error {
	var castedVal any
	var err error

	if kind == "bool" {
		if value == "" {
			value = "true"
		}

		castedVal, err = strconv.ParseBool(value)

		if err != nil {
			return err
		}
	} else if kind == "int64" {
		if value == "" {
			value = "0"
		}

		castedVal, err = strconv.ParseInt(value, 10, 64)

		if err != nil {
			return err
		}
	} else if kind == "uint64" {
		if value == "" {
			value = "0"
		}

		castedVal, err = strconv.ParseUint(value, 10, 64)

		if err != nil {
			return err
		}
	} else if kind == "slice" {
		castedVal = strings.Split(value, ",")
	} else if kind == "string" {
		castedVal = value
	} else {
		return nil
	}

	if !funk.ContainsString(fc.availableSpecs, specName) {
		fc.availableSpecs = append(fc.availableSpecs, specName)
	}

	return attr.SetValue(fc, specName, castedVal)
}

func (fc *FieldSpecs) Parse(model any, field string, kinds map[string]string) error {
	var err error

	fc.Name = field

	fc.BSON, _ = attr.GetTag(model, field, "bson")
	fc.JSON, _ = attr.GetTag(model, field, "json")

	// extract field name (for exmaple _id) from
	// _id,omitempty and similar values
	fc.BSON = utils.StringUtilsInstance.StringFirstToken(fc.BSON, ",")
	fc.JSON = utils.StringUtilsInstance.StringFirstToken(fc.JSON, ",")

	fc.BSON = strings.TrimSpace(fc.BSON)
	fc.JSON = strings.TrimSpace(fc.JSON)

	// Iterate all kinds of the fields
	// and set specs by reading it from the model's field
	for iname, ikind := range kinds {
		inameLower := strings.ToLower(iname)
		inameLower = fc.TranslateField(inameLower)

		fullName := "go-rh-" + inameLower

		val, _ := attr.GetTag(model, field, fullName)

		if val != "" {
			if err = fc.setSpecByKind(iname, val, ikind); err != nil {
				return err
			}
		}
	}

	// Iterate model's specs in go-rh tag
	// and set the values
	val, _ := attr.GetTag(model, field, "go-rh")

	for _, isubVal := range strings.Split(val, ";") {
		// 0 name
		// 1 value
		isubValParts := strings.SplitN(isubVal, ":", 2)

		isubValParts[0] = strings.TrimSpace(isubValParts[0])

		for iname, ikind := range kinds {
			inameLower := strings.ToLower(iname)
			inameLower = fc.TranslateField(inameLower)

			if inameLower != isubValParts[0] {
				continue
			}

			val := ""

			if len(isubValParts) == 2 {
				val = isubValParts[1]
			}

			if err = fc.setSpecByKind(iname, val, ikind); err != nil {
				return err
			}
		}
	}

	value, err := attr.GetValue(model, field)

	if err != nil {
		return err
	}

	fc.fieldType = reflect.TypeOf(value)

	return nil
}
