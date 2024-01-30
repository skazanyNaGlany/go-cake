package go_cake

import (
	"fmt"
	"log"
	"regexp"

	"github.com/skazanyNaGlany/go-cake/utils"
	"github.com/thoas/go-funk"
)

type Resource struct {
	Pattern, DbPath, ResourceName string
	DatabaseDriver                DatabaseDriver
	DbModel                       GoKateModel
	DbModelIDField                string // TODO test if field exists in the model
	DbModelETagField              string // TODO test if field exists in the model
	SupportedVersion              []string
	DbModelJSONFields             []string
	DbModelJSONFieldsNoReserved   []string
	CompiledPattern               *regexp.Regexp
	UserData                      any
	ResourceCallback              *ResourceCallback
	JSONSchemaConfig              *JSONSchemaConfig
	CORSConfig                    *CORSConfig
	GetAllowed                    bool
	DeleteAllowed                 bool
	InsertAllowed                 bool
	UpdateAllowed                 bool
	GetMaxOutputItems             int64
	DeleteMaxInputItems           int64
	DeleteMaxInputPayloadSize     int64
	InsertMaxInputItems           int64
	InsertMaxInputPayloadSize     int64
	UpdateMaxInputItems           int64
	UpdateMaxInputPayloadSize     int64
	compiledSupportedVersion      []*regexp.Regexp
}

func (rhr *Resource) Init() error {
	var err error

	if err = rhr.testResource(); err != nil {
		return err
	}

	rhr.CompiledPattern, err = regexp.Compile(rhr.Pattern)

	if err != nil {
		return err
	}

	if err = rhr.checkSchemaConfig(); err != nil {
		return err
	}

	if err = rhr.collectJSONFields(); err != nil {
		return err
	}

	if err = rhr.testResource2(); err != nil {
		return err
	}

	if err = rhr.checkSchemaConfigFields(); err != nil {
		return err
	}

	if err = rhr.initJSONValidator(); err != nil {
		return err
	}

	if err = rhr.compileSupportedVersions(); err != nil {
		return err
	}

	rhr.initRanges()

	if err = rhr.testModel(); err != nil {
		return err
	}

	return nil
}

func (rhr *Resource) testResource() error {
	if rhr.Pattern == "" {
		return NewNoResourcePatternSetError(rhr, nil)
	}

	if rhr.DbPath == "" {
		return NewNoResourceDbPathSetError(rhr, nil)
	}

	if rhr.ResourceName == "" {
		return NewNoResourceResourceNameSetError(rhr, nil)
	}

	if rhr.DatabaseDriver == nil {
		return NewNoResourceDatabaseDriverSetError(rhr, nil)
	}

	if rhr.DbModel == nil {
		return NewNoResourceDbModelSetError(rhr, nil)
	}

	if rhr.JSONSchemaConfig == nil {
		return NewNoResourceJSONSchemaConfigSetError(rhr, nil)
	}

	if rhr.JSONSchemaConfig.IDField == "" {
		return NewNoResourceJSONSchemaConfigIDFieldSetError(rhr, nil)
	}

	return nil
}

func (rhr *Resource) testModel() error {
	if err := rhr.DatabaseDriver.TestModel(
		rhr.DbModelIDField,
		rhr.DbModelETagField,
		rhr.DbModel,
		rhr.DbPath); err != nil {
		return NewUnableToTestModelError(rhr, rhr.DatabaseDriver, rhr.DbModel, err)
	}

	return nil
}

func (rhr *Resource) MatchPattern(path string) bool {
	return rhr.CompiledPattern.MatchString(path)
}

func (rhr *Resource) collectJSONFields() error {
	var err error

	tagMap, err := utils.StructUtilsInstance.StructToTagMap(rhr.DbModel, []string{"json"}, "name")

	if err != nil {
		return err
	}

	jsonIdField := rhr.JSONSchemaConfig.IDField
	jsonEtagField := rhr.JSONSchemaConfig.ETagField

	for _, tags := range tagMap {
		fieldNameInJsonTag, hasJsonTag := tags["json"]

		if !hasJsonTag {
			continue
		}

		if fieldNameInJsonTag == "-" {
			continue
		}

		rhr.DbModelJSONFields = append(rhr.DbModelJSONFields, fieldNameInJsonTag)

		if fieldNameInJsonTag != jsonIdField && fieldNameInJsonTag != jsonEtagField {
			rhr.DbModelJSONFieldsNoReserved = append(rhr.DbModelJSONFieldsNoReserved, fieldNameInJsonTag)
		}
	}

	return nil
}

func (rhr *Resource) testResource2() error {
	jsonIdField := rhr.JSONSchemaConfig.IDField
	jsonEtagField := rhr.JSONSchemaConfig.ETagField

	if !funk.ContainsString(rhr.DbModelJSONFields, jsonIdField) {
		return NewIDFieldModelNotFoundError(rhr, rhr.DbModel, jsonIdField, nil)
	}

	if jsonEtagField != "" {
		if !funk.ContainsString(rhr.DbModelJSONFields, jsonEtagField) {
			return NewETagFieldModelNotFoundError(rhr, rhr.DbModel, jsonEtagField, nil)
		}
	}

	return nil
}

func (rhr *Resource) checkSchemaConfig() error {
	if rhr.JSONSchemaConfig == nil {
		return NoSchemaConfigError{}
	}

	if rhr.JSONSchemaConfig.IDField == "" {
		return NoSchemaConfigIDError{}
	}

	return nil
}

func (rhr *Resource) initJSONValidator() error {
	if rhr.JSONSchemaConfig == nil {
		return nil
	}

	if rhr.JSONSchemaConfig.Validator == nil {
		return nil
	}

	return nil
}

func (rhr *Resource) compileSupportedVersions() error {
	for _, pattern := range rhr.SupportedVersion {
		compiled, err := regexp.Compile(pattern)

		if err != nil {
			return err
		}

		rhr.compiledSupportedVersion = append(rhr.compiledSupportedVersion, compiled)
	}

	return nil
}

func (rhr *Resource) checkSchemaConfigFields() error {
	allFields := rhr.JSONSchemaConfig.GetAllFields()

	// log.Println("allFields", allFields)
	// log.Println("rhr.DbModelJSONFields", rhr.DbModelJSONFields)

	for _, iField := range allFields {
		if iField == FIELD_ANY {
			continue
		}

		if !funk.ContainsString(rhr.DbModelJSONFields, iField) {
			log.Println("allFields", allFields)
			log.Println("rhr.DbModelJSONFields", rhr.DbModelJSONFields)
			log.Println("iField", iField)

			return SchemaConfigUnknownFieldError{
				BaseError{Message: fmt.Sprintf("Unknown field %v in %T model", iField, rhr.DbModel)}}
		}
	}

	return nil
}

func (rhr *Resource) initRanges() {
	if rhr.GetMaxOutputItems == 0 {
		rhr.GetMaxOutputItems = MAX_OUTPUT_ITEMS
	}

	if rhr.DeleteMaxInputItems == 0 {
		rhr.DeleteMaxInputItems = MAX_INPUT_ITEMS
	}

	if rhr.DeleteMaxInputPayloadSize == 0 {
		rhr.DeleteMaxInputPayloadSize = MAX_INPUT_PAYLOAD_SIZE
	}

	if rhr.InsertMaxInputItems == 0 {
		rhr.InsertMaxInputItems = MAX_INPUT_ITEMS
	}

	if rhr.InsertMaxInputPayloadSize == 0 {
		rhr.InsertMaxInputPayloadSize = MAX_INPUT_PAYLOAD_SIZE
	}

	if rhr.UpdateMaxInputItems == 0 {
		rhr.UpdateMaxInputItems = MAX_INPUT_ITEMS
	}

	if rhr.UpdateMaxInputPayloadSize == 0 {
		rhr.UpdateMaxInputPayloadSize = MAX_INPUT_PAYLOAD_SIZE
	}
}
