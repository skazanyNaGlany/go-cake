package go_cake

import "net/http"

// app handlers
type AuthAppFunc func(
	resource *Resource,
	request *Request,
	response *ResponseJSON) bool

type PrePostRequestAppFunc func(
	resource *Resource,
	request *Request,
	response *ResponseJSON) HTTPError

type DocumentsAppFunc func(
	resource *Resource,
	request *Request,
	documents []GoKateModel,
	currentHttpErr HTTPError) HTTPError

// internal handlers
type MiddlewareFunc func(http.Handler) http.Handler
