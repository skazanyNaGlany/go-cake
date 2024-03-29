package go_cake

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/skazanyNaGlany/go-cake/utils"
	"github.com/thoas/go-funk"
)

type BaseRequestProcessor struct {
	request             *Request
	resource            *Resource
	subRequestProcessor RequestProcessor
}

func (brp *BaseRequestProcessor) ProcessRequest(response *ResponseJSON) {
	var httpErr HTTPError = NewOKHTTPError(nil)
	var documents []GoCakeModel

	timeStart := time.Now()

	response.Meta.RequestUniqueID = brp.request.UniqueID
	response.Meta.Version = brp.request.Version
	response.Meta.URL = brp.request.URL
	response.Meta.Method = brp.request.Method
	response.Meta.StatusMessage = httpErr.GetStatusMessage()
	response.Meta.StatusCode = httpErr.GetStatusCode()

	defer func() {
		defer func() {
			brp.catchInternalError(response, recover())
		}()

		if brp.catchInternalError(response, recover()) {
			return
		}

		brp.postRequestResponseActions(response)

		response.Meta.TotalTimeMs = time.Since(timeStart).Seconds() * 1000

		brp.processTotals(response)

		if httpErr = brp.callPostRequestHandlers(response); httpErr != nil {
			response.Meta.StatusMessage = httpErr.GetStatusMessage()
			response.Meta.StatusCode = httpErr.GetStatusCode()
		}
	}()

	if httpErr = brp.processCORS(); httpErr != nil {
		response.Meta.StatusMessage = httpErr.GetStatusMessage()
		response.Meta.StatusCode = httpErr.GetStatusCode()

		return
	}

	if httpErr = brp.checkSupportedVersion(); httpErr != nil {
		response.Meta.StatusMessage = httpErr.GetStatusMessage()
		response.Meta.StatusCode = httpErr.GetStatusCode()

		return
	}

	if httpErr = brp.preRequestProjectableChecks(); httpErr != nil {
		response.Meta.StatusMessage = httpErr.GetStatusMessage()
		response.Meta.StatusCode = httpErr.GetStatusCode()

		return
	}

	if httpErr = brp.preRequestFilterableChecks(); httpErr != nil {
		response.Meta.StatusMessage = httpErr.GetStatusMessage()
		response.Meta.StatusCode = httpErr.GetStatusCode()

		return
	}

	if httpErr = brp.preRequestSortableChecks(); httpErr != nil {
		response.Meta.StatusMessage = httpErr.GetStatusMessage()
		response.Meta.StatusCode = httpErr.GetStatusCode()

		return
	}

	if httpErr = brp.callPreRequestHandlers(response); httpErr != nil {
		response.Meta.StatusMessage = httpErr.GetStatusMessage()
		response.Meta.StatusCode = httpErr.GetStatusCode()

		return
	}

	if httpErr = brp.callAuthHandlers(response); httpErr != nil {
		response.Meta.StatusMessage = httpErr.GetStatusMessage()
		response.Meta.StatusCode = httpErr.GetStatusCode()

		return
	}

	if documents, httpErr = brp.subRequestProcessor.ProcessRequest(response); httpErr != nil {
		response.Meta.StatusMessage = httpErr.GetStatusMessage()
		response.Meta.StatusCode = httpErr.GetStatusCode()
	}

	brp.documentsToJsonMapObjects(documents, response)
}

func (brp *BaseRequestProcessor) catchInternalError(response *ResponseJSON, r any) bool {
	if r != nil {
		message := fmt.Sprint(r)

		httpErr := NewInternalServerErrorHTTPError(errors.New(message))

		response.Meta.StatusMessage = httpErr.GetStatusMessage()
		response.Meta.StatusCode = httpErr.GetStatusCode()

		return true
	}

	return false
}

func (brp *BaseRequestProcessor) processTotals(response *ResponseJSON) {
	ctx, cancel := brp.resource.ResourceCallback.CreateContext(
		brp.resource,
		brp.request,
		response,
		ctxDbDriverTotal)
	defer cancel()

	response.Meta.Total, _ = brp.resource.DatabaseDriver.Total(
		brp.resource.DbModel,
		brp.request.Where,
		ctx,
		nil)
}

func (brp *BaseRequestProcessor) checkSupportedVersion() HTTPError {
	if !utils.RegExUtilsInstance.HasMatch(
		brp.resource.compiledSupportedVersion,
		brp.request.Version) {
		return NewUnsupportedVersionHTTPError(brp.request.Version, nil)
	}

	return nil
}

func (brp *BaseRequestProcessor) processCORS() HTTPError {
	if brp.resource.CORSConfig == nil {
		return nil
	}

	if brp.request.IsCORS {
		return brp.processCORSbyCORSRequest()
	} else {
		return brp.processCORSbyOtherRequest()
	}
}

func (brp *BaseRequestProcessor) processCORSbyCORSRequest() HTTPError {
	origin := brp.request.Request.Header.Get("Origin")
	supportedMethods := brp.getSupportedMethodsForOrigin(origin)

	brp.CORSwriteAccessControlAllowMethods(
		brp.request.ResponseWriter,
		supportedMethods)

	method := brp.request.Request.Header.Get("Access-Control-Request-Method")
	methodIsSupported := funk.ContainsString(supportedMethods, method)

	if !methodIsSupported {
		brp.CORSwriteAccessControlAllowOrigin(brp.request.ResponseWriter, "null")
	} else {
		brp.CORSwriteAccessControlAllowOrigin(brp.request.ResponseWriter, origin)
	}

	return nil
}

func (brp *BaseRequestProcessor) processCORSbyOtherRequest() HTTPError {
	origin := brp.request.Request.Header.Get("Origin")
	supportedMethods := brp.getSupportedMethodsForOrigin(origin)

	brp.CORSwriteAccessControlAllowMethods(
		brp.request.ResponseWriter,
		supportedMethods)

	method := brp.request.Method
	methodIsSupported := funk.ContainsString(supportedMethods, method)

	if !methodIsSupported {
		brp.CORSwriteAccessControlAllowOrigin(brp.request.ResponseWriter, "null")
	} else {
		brp.CORSwriteAccessControlAllowOrigin(brp.request.ResponseWriter, origin)
	}

	if !methodIsSupported {
		return NewMethodNotAllowedHTTPError(nil)
	}

	return nil
}

func (brp *BaseRequestProcessor) callAuthHandlers(response *ResponseJSON) HTTPError {
	if brp.resource.ResourceCallback == nil ||
		brp.resource.ResourceCallback.AuthCallback == nil {
		return nil
	}

	if !brp.resource.ResourceCallback.AuthCallback(
		brp.resource,
		brp.request,
		response) {
		return NewUnauthorizedHTTPError(nil)
	}

	return nil
}

func (brp *BaseRequestProcessor) callFetchedDocumentsHandlers(
	documents []GoCakeModel,
	currentHttpErr HTTPError) HTTPError {
	if brp.resource.ResourceCallback == nil ||
		brp.resource.ResourceCallback.FetchedDocuments == nil {
		return nil
	}

	return brp.resource.ResourceCallback.FetchedDocuments(
		brp.resource,
		brp.request,
		documents,
		currentHttpErr)
}

func (brp *BaseRequestProcessor) callUpdatingDocumentsHandlers(
	documents []GoCakeModel,
	currentHttpErr HTTPError) HTTPError {
	if brp.resource.ResourceCallback == nil ||
		brp.resource.ResourceCallback.UpdatingDocuments == nil {
		return nil
	}

	return brp.resource.ResourceCallback.UpdatingDocuments(
		brp.resource,
		brp.request,
		documents,
		currentHttpErr)
}

func (brp *BaseRequestProcessor) callUpdatedDocumentsHandlers(
	documents []GoCakeModel,
	currentHttpErr HTTPError) HTTPError {
	if brp.resource.ResourceCallback == nil ||
		brp.resource.ResourceCallback.UpdatedDocuments == nil {
		return nil
	}

	return brp.resource.ResourceCallback.UpdatedDocuments(
		brp.resource,
		brp.request,
		documents,
		currentHttpErr)
}

func (brp *BaseRequestProcessor) callInsertingDocumentsHandlers(
	documents []GoCakeModel, currentHttpErr HTTPError) HTTPError {
	if brp.resource.ResourceCallback == nil ||
		brp.resource.ResourceCallback.InsertingDocuments == nil {
		return nil
	}

	return brp.resource.ResourceCallback.InsertingDocuments(
		brp.resource,
		brp.request,
		documents,
		currentHttpErr)
}

func (brp *BaseRequestProcessor) callInsertedDocumentsHandlers(
	documents []GoCakeModel,
	currentHttpErr HTTPError) HTTPError {
	if brp.resource.ResourceCallback == nil ||
		brp.resource.ResourceCallback.InsertedDocuments == nil {
		return nil
	}

	return brp.resource.ResourceCallback.InsertedDocuments(
		brp.resource,
		brp.request,
		documents,
		currentHttpErr)
}

func (brp *BaseRequestProcessor) callDeletingDocumentsHandlers(
	documents []GoCakeModel,
	currentHttpErr HTTPError) HTTPError {
	if brp.resource.ResourceCallback == nil ||
		brp.resource.ResourceCallback.DeletingDocuments == nil {
		return nil
	}

	return brp.resource.ResourceCallback.DeletingDocuments(
		brp.resource,
		brp.request,
		documents,
		currentHttpErr)
}

func (brp *BaseRequestProcessor) callDeletedDocumentsHandlers(
	documents []GoCakeModel,
	currentHttpErr HTTPError) HTTPError {
	if brp.resource.ResourceCallback == nil ||
		brp.resource.ResourceCallback.DeletedDocuments == nil {
		return nil
	}

	return brp.resource.ResourceCallback.DeletedDocuments(
		brp.resource,
		brp.request,
		documents,
		currentHttpErr)
}

func (brp *BaseRequestProcessor) callPreRequestHandlers(response *ResponseJSON) HTTPError {
	if brp.resource.ResourceCallback == nil ||
		brp.resource.ResourceCallback.PreRequestCallback == nil {
		return nil
	}

	return brp.resource.ResourceCallback.PreRequestCallback(
		brp.resource,
		brp.request,
		response)
}

func (brp *BaseRequestProcessor) callPostRequestHandlers(response *ResponseJSON) HTTPError {
	if brp.resource.ResourceCallback == nil ||
		brp.resource.ResourceCallback.PostRequestCallback == nil {
		return nil
	}

	return brp.resource.ResourceCallback.PostRequestCallback(
		brp.resource,
		brp.request,
		response)
}

func (brp *BaseRequestProcessor) findNonExistingFields(
	fields []string,
	s []string) []string {
	nonExistingFields := make([]string, 0)

	for _, iField := range fields {
		if !slices.Contains(s, iField) {
			nonExistingFields = append(nonExistingFields, iField)
		}
	}

	return nonExistingFields
}

func (brp *BaseRequestProcessor) preRequestFilterableChecks() HTTPError {
	if brp.request.Where == "" {
		return nil
	}

	whereFields, httpErr := brp.resource.DatabaseDriver.GetWhereFields(
		brp.resource.DbModel,
		brp.request.Where)

	if httpErr != nil {
		return httpErr
	}

	nonExistingFields := brp.findNonExistingFields(
		whereFields,
		brp.resource.DbModelJSONFields)

	if len(nonExistingFields) > 0 {
		return NewFieldNotExistsHTTPError(nonExistingFields[0], nil)
	}

	filterableFields := brp.resource.JSONSchemaConfig.FilterableFields

	if funk.ContainsString(filterableFields, FIELD_ANY) {
		return nil
	}

	for _, iJsonField := range whereFields {
		if !funk.ContainsString(filterableFields, iJsonField) {
			return NewFieldNotFilterableHTTPError(iJsonField, nil)
		}
	}

	return nil
}

func (brp *BaseRequestProcessor) preRequestSortableChecks() HTTPError {
	if brp.request.Sort == "" {
		return nil
	}

	sortFields, httpErr := brp.resource.DatabaseDriver.GetSortFields(
		brp.resource.DbModel,
		brp.request.Sort)

	if httpErr != nil {
		return httpErr
	}

	nonExistingFields := brp.findNonExistingFields(
		sortFields,
		brp.resource.DbModelJSONFields)

	if len(nonExistingFields) > 0 {
		return NewFieldNotExistsHTTPError(nonExistingFields[0], nil)
	}

	sortableFields := brp.resource.JSONSchemaConfig.SortableFields

	if funk.ContainsString(sortableFields, FIELD_ANY) {
		return nil
	}

	for _, iJsonField := range sortFields {
		if !funk.ContainsString(sortableFields, iJsonField) {
			return NewFieldNotSortableHTTPError(iJsonField, nil)
		}
	}

	return nil
}

func (brp *BaseRequestProcessor) preRequestProjectableChecks() HTTPError {
	if len(brp.request.Projection) == 0 {
		return nil
	}

	nonExistingFields := brp.findNonExistingFields(
		brp.request.ProjectionFields,
		brp.resource.DbModelJSONFields)

	if len(nonExistingFields) > 0 {
		return NewFieldNotExistsHTTPError(nonExistingFields[0], nil)
	}

	projectableFields := brp.resource.JSONSchemaConfig.ProjectableFields

	if funk.ContainsString(projectableFields, FIELD_ANY) {
		return nil
	}

	for iJsonField := range brp.request.Projection {
		if !funk.ContainsString(projectableFields, iJsonField) {
			return NewFieldNotProjectableHTTPError(iJsonField, nil)
		}
	}

	return nil
}

func (brp *BaseRequestProcessor) preRequestRequireOnInsertChecks(jsonObjectMap map[string]any, requiredFields []string) HTTPError {
	for _, irequiredField := range requiredFields {
		if _, ok := jsonObjectMap[irequiredField]; !ok {
			return NewFieldRequiredHTTPError(irequiredField, nil)
		}
	}

	return nil
}

func (brp *BaseRequestProcessor) preRequestRequireOnUpdateChecks(
	jsonObjectMap map[string]any,
	requiredFields []string) HTTPError {
	for _, irequiredField := range requiredFields {
		if _, ok := jsonObjectMap[irequiredField]; !ok {
			return NewFieldRequiredHTTPError(irequiredField, nil)
		}
	}

	return nil
}

func (brp *BaseRequestProcessor) preRequestOptimizeFields(
	jsonObjectMap map[string]any,
	optimizeFields []string,
	anyField bool) {
	for iJsonField, iJsonValue := range jsonObjectMap {
		if iJsonValue == nil {
			continue
		}

		if !anyField {
			if !funk.ContainsString(optimizeFields, iJsonField) {
				continue
			}
		}

		fieldType := reflect.TypeOf(iJsonValue)

		if fieldType.Kind() == reflect.String {
			jsonObjectMap[iJsonField] = utils.StringUtilsInstance.OptimizeString(iJsonValue.(string))
		}
	}
}

func (brp *BaseRequestProcessor) preRequestInsertableChecks(jsonObjectMap map[string]any, requiredFields, insertableFields []string) HTTPError {
	for iJsonField := range jsonObjectMap {
		if funk.ContainsString(requiredFields, iJsonField) {
			// field is required, so do not check if it is insertable or not
			continue
		}

		if !funk.ContainsString(insertableFields, iJsonField) {
			return NewFieldNotInsertableHTTPError(iJsonField, nil)
		}
	}

	return nil
}

func (brp *BaseRequestProcessor) preRequestUpdatableChecks(
	jsonObjectMap map[string]any,
	requiredFields,
	updatableFields []string) HTTPError {
	jsonIdField := brp.resource.JSONSchemaConfig.IDField
	jsonEtagField := brp.resource.JSONSchemaConfig.ETagField

	for iJsonField := range jsonObjectMap {
		if iJsonField == jsonIdField || iJsonField == jsonEtagField {
			// allow ID/ETAG field in the payload
			// it will be removed at the driver eventually
			continue
		}

		if funk.ContainsString(requiredFields, iJsonField) {
			// field is required, so do not check if it is updatable or not
			continue
		}

		if !funk.ContainsString(updatableFields, iJsonField) {
			return NewFieldNotUpdatableHTTPError(iJsonField, nil)
		}
	}

	return nil
}

func (brp *BaseRequestProcessor) documentsToJsonMapObjects(documents []GoCakeModel, response *ResponseJSON) HTTPError {
	for _, iDoc := range documents {
		jsonObjectMap, _ := iDoc.ToMap()

		response.Items = append(response.Items, jsonObjectMap)
	}

	return nil
}

func (brp *BaseRequestProcessor) checkDocumentsForErrors(documents []GoCakeModel) HTTPError {
	for _, iDocument := range documents {
		if iDocument.GetHTTPError() != nil {
			return NewPayloadInvalidHTTPError(nil)
		}
	}

	return nil
}

func (brp *BaseRequestProcessor) postRequestResponseHiddenAction(jsonObjectMap map[string]any, projecionFields map[string]bool, hiddenFields []string) {
	for _, iHiddenField := range hiddenFields {
		if projected, keyIn := projecionFields[iHiddenField]; keyIn {
			if projected {
				// field is projected at the URL
				continue
			}
		}

		if _, keyIn := jsonObjectMap[iHiddenField]; !keyIn {
			continue
		}

		delete(jsonObjectMap, iHiddenField)
	}
}

func (brp *BaseRequestProcessor) postRequestResponseProjectableAction(jsonObjectMap map[string]any, projecionFields map[string]bool) {
	for iProjectedField, iProjectedFieldValue := range projecionFields {
		if _, keyIn := jsonObjectMap[iProjectedField]; !keyIn {
			continue
		}

		if iProjectedFieldValue {
			continue
		}

		delete(jsonObjectMap, iProjectedField)
	}
}

func (brp *BaseRequestProcessor) postRequestResponseErasedAction(jsonObjectMap map[string]any, erasedFields []string) {
	for _, iErasedField := range erasedFields {
		if _, keyIn := jsonObjectMap[iErasedField]; !keyIn {
			continue
		}

		if _, isStr := jsonObjectMap[iErasedField].(string); !isStr {
			// value is not a string
			continue
		}

		jsonObjectMap[iErasedField] = ""
	}
}

func (brp *BaseRequestProcessor) postRequestResponseActions(response *ResponseJSON) HTTPError {
	hiddenFields := brp.resource.JSONSchemaConfig.HiddenFields
	erasedFields := brp.resource.JSONSchemaConfig.ErasedFields

	hiddenAnyField := funk.ContainsString(hiddenFields, FIELD_ANY)
	erasedAnyField := funk.ContainsString(erasedFields, FIELD_ANY)

	if hiddenAnyField {
		hiddenFields = brp.resource.DbModelJSONFieldsNoReserved
	}

	if erasedAnyField {
		erasedFields = brp.resource.DbModelJSONFieldsNoReserved
	}

	// projection fields was validated at preRequestProjectableChecks()
	for _, jsonObject := range response.Items {
		brp.postRequestResponseHiddenAction(
			jsonObject,
			brp.request.Projection,
			hiddenFields)

		brp.postRequestResponseProjectableAction(
			jsonObject, brp.request.Projection)

		brp.postRequestResponseErasedAction(
			jsonObject,
			erasedFields)
	}

	return nil
}

func (brp *BaseRequestProcessor) CORSwriteAccessControlAllowOrigin(
	responseWriter http.ResponseWriter,
	origin string) {
	responseWriter.Header().Set("Access-Control-Allow-Origin", origin)
}

func (brp *BaseRequestProcessor) CORSwriteAccessControlAllowMethods(
	responseWriter http.ResponseWriter,
	supportedMethods []string) {

	responseWriter.Header().Set(
		"Access-Control-Allow-Methods",
		strings.Join(supportedMethods, ", "))
}

func (brp *BaseRequestProcessor) getSupportedMethodsForOrigin(origin string) []string {
	supported := make([]string, 0)
	supported = append(supported, brp.request.GetCORSMethods()...)

	// Get
	if utils.RegExUtilsInstance.HasMatch(
		brp.resource.CORSConfig.getCompiledOrigins,
		origin) {
		supported = append(supported, brp.request.GetGetMethods()...)
	}

	// Delete
	if utils.RegExUtilsInstance.HasMatch(
		brp.resource.CORSConfig.deleteCompiledOrigins,
		origin) {
		supported = append(supported, brp.request.GetDeleteMethods()...)
	}

	// Insert
	if utils.RegExUtilsInstance.HasMatch(
		brp.resource.CORSConfig.insertCompiledOrigins,
		origin) {
		supported = append(supported, brp.request.GetInsertMethods()...)
	}

	// Update
	if utils.RegExUtilsInstance.HasMatch(
		brp.resource.CORSConfig.updateCompiledOrigins,
		origin) {
		supported = append(supported, brp.request.GetUpdateMethods()...)
	}

	return supported
}

func (brp *BaseRequestProcessor) decodedJsonSliceToDBModels() ([]GoCakeModel, HTTPError) {
	converted := make([]GoCakeModel, 0)

	for _, jsonObject := range brp.request.DecodedJsonSlice {
		modelNewInstance := brp.resource.DbModel.CreateInstance()

		jsonObjectBytes, err := json.Marshal(jsonObject)

		if err != nil {
			modelNewInstance.SetHTTPError(
				NewClientObjectMalformedHTTPError(err))
		}

		if modelNewInstance.GetHTTPError() == nil {
			if err = json.Unmarshal(jsonObjectBytes, modelNewInstance); err != nil {
				modelNewInstance.SetHTTPError(
					NewClientObjectMalformedHTTPError(err))
			}
		}

		if modelNewInstance.GetHTTPError() == nil {
			if httpErr := jsonObject["__http_error__"]; httpErr != nil {
				modelNewInstance.SetHTTPError(httpErr.(HTTPError))
			}
		}

		converted = append(converted, modelNewInstance)
	}

	return converted, nil
}
