package go_cake

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/skazanyNaGlany/go-cake/utils"
	"github.com/thoas/go-funk"
)

var allowedAcceptValues []string = []string{
	ALLOWED_ACCEPT_HEADER_0,
	ALLOWED_ACCEPT_HEADER_1}

type Request struct {
	ResourcePattern  *regexp.Regexp
	Version          string
	Resource         string
	Where            string
	Sort             string
	Projection       string
	Page             int64
	PerPage          int64
	UniqueID         string
	Method           string
	URL              string
	Body             []byte
	DecodedJsonSlice []map[string]any
	ContentLength    int64
	Request          *http.Request
	ResponseWriter   http.ResponseWriter
	UserData         any
	IsGet            bool
	IsInsert         bool
	IsUpdate         bool
	IsDelete         bool
	IsCORS           bool
}

func (rhr Request) HasWhere() bool {
	return rhr.Where != ""
}

func (rhr Request) HasSort() bool {
	return rhr.Sort != ""
}

func (rhr Request) HasPage() bool {
	return rhr.Page > 0
}

func (rhr *Request) Parse(r *http.Request) HTTPError {
	var err error
	var httpErr HTTPError

	rhr.UniqueID = utils.StringUtilsInstance.NewUUID()

	if len(r.URL.Path) > MAX_URL_LENGTH {
		return NewURLTooBigHTTPError(MAX_URL_LENGTH, nil)
	}

	if httpErr = rhr.checkAcceptHeader(r); httpErr != nil {
		return httpErr
	}

	rhr.Method = r.Method
	rhr.URL = r.URL.String()
	rhr.ContentLength = r.ContentLength

	if rhr.MethodIsGet(rhr.Method) {
		rhr.IsGet = true
	} else if rhr.MethodIsDelete(rhr.Method) {
		rhr.IsDelete = true
	} else if rhr.MethodIsInsert(rhr.Method) {
		rhr.IsInsert = true
	} else if rhr.MethodIsUpdate(rhr.Method) {
		rhr.IsUpdate = true
	} else if rhr.MethodIsCORS(rhr.Method) {
		rhr.IsCORS = true
	}

	if rhr.IsDelete || rhr.IsInsert || rhr.IsUpdate {
		if httpErr := rhr.checkContentTypeHeader(r); httpErr != nil {
			return httpErr
		}
	}

	urlParts := utils.RegExUtilsInstance.FindNamedMatches(
		rhr.ResourcePattern,
		r.URL.Path)

	if _, ok := urlParts["url"]; !ok {
		return NewUnableToParseRequestHTTPError(nil)
	}

	if version, ok := urlParts["version"]; ok {
		version = strings.Trim(version, "/")

		rhr.Version = strings.TrimSpace(version)
	}

	query := r.URL.Query()

	where := strings.TrimSpace(query.Get("where"))
	sort := strings.TrimSpace(query.Get("sort"))
	projection := strings.TrimSpace(query.Get("projection"))
	perPage := strings.TrimSpace(query.Get("per_page"))
	page := strings.TrimSpace(query.Get("page"))

	if where != "" {
		rhr.Where = where
	}

	if sort != "" {
		rhr.Sort = sort
	}

	if projection != "" {
		rhr.Projection = projection
	}

	if page != "" {
		rhr.Page, _ = strconv.ParseInt(page, 10, 64)
	}

	if perPage != "" {
		rhr.PerPage, _ = strconv.ParseInt(perPage, 10, 64)
	}

	rhr.Body, err = io.ReadAll(r.Body)

	r.Body.Close()

	if err != nil {
		return NewUnableToParseRequestHTTPError(err)
	}

	rhr.DecodedJsonSlice, httpErr = rhr.requestBodyToArrayOfMaps()

	if httpErr != nil {
		return httpErr
	}

	return nil
}

func (rhr *Request) GetGetMethods() []string {
	return []string{HTTP_REQUEST_GET_METHOD}
}

func (rhr *Request) GetDeleteMethods() []string {
	return []string{HTTP_REQUEST_DELETE_METHOD}
}

func (rhr *Request) GetInsertMethods() []string {
	return []string{HTTP_REQUEST_POST_METHOD, HTTP_REQUEST_PUT_METHOD}
}

func (rhr *Request) GetUpdateMethods() []string {
	return []string{HTTP_REQUEST_PATCH_METHOD}
}

func (rhr *Request) GetCORSMethods() []string {
	return []string{HTTP_REQUEST_OPTIONS_METHOD}
}

func (rhr *Request) MethodIsGet(method string) bool {
	return funk.ContainsString(rhr.GetGetMethods(), method)
}

func (rhr *Request) MethodIsDelete(method string) bool {
	return funk.ContainsString(rhr.GetDeleteMethods(), method)
}

func (rhr *Request) MethodIsInsert(method string) bool {
	return funk.ContainsString(rhr.GetInsertMethods(), method)
}

func (rhr *Request) MethodIsUpdate(method string) bool {
	return funk.ContainsString(rhr.GetUpdateMethods(), method)
}

func (rhr *Request) MethodIsCORS(method string) bool {
	return funk.ContainsString(rhr.GetCORSMethods(), method)
}

func (rhr *Request) checkAcceptHeader(r *http.Request) HTTPError {
	values, ok := r.Header["Accept"]

	if !ok {
		return NewInvalidAcceptRequestHeaderHTTPError(allowedAcceptValues, nil)
	}

	for _, ivalue := range values {
		for _, iSubValue := range strings.Split(ivalue, ",") {
			iSubValue = strings.TrimSpace(iSubValue)
			iSubValueParts := strings.Split(iSubValue, ";")

			if len(iSubValueParts) > 0 {
				iSubValueParts[0] = strings.TrimSpace(iSubValueParts[0])

				if funk.ContainsString(allowedAcceptValues, iSubValueParts[0]) {
					return nil
				}
			}
		}
	}

	return NewInvalidAcceptRequestHeaderHTTPError(allowedAcceptValues, nil)
}

func (rhr *Request) checkContentTypeHeader(r *http.Request) HTTPError {
	values := r.Header["Content-Type"]

	if !funk.ContainsString(values, ALLOWED_REQUEST_CONTENT_TYPE) {
		return NewInvalidContentTypeRequestHeaderHTTPError(ALLOWED_REQUEST_CONTENT_TYPE, nil)
	}

	return nil
}

func (rhr *Request) requestBodyToArrayOfMaps() ([]map[string]any, HTTPError) {
	var decodedSlice []map[string]any
	var decodedObject map[string]any

	if len(rhr.Body) == 0 {
		return decodedSlice, nil
	}

	if err := json.Unmarshal(rhr.Body, &decodedSlice); err != nil {
		if err := json.Unmarshal(rhr.Body, &decodedObject); err != nil {
			httpErr := NewCannotDecodePayloadHTTPError(err)

			return decodedSlice, httpErr
		}
	}

	if len(decodedSlice) == 0 && reflect.ValueOf(decodedObject).Kind() == reflect.Map {
		decodedSlice = append(decodedSlice, decodedObject)
	}

	return decodedSlice, nil
}
