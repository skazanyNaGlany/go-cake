package go_cake

import "net/http"

type MiddlewareQueue struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Queue          []MiddlewareCallback
}

func (mq *MiddlewareQueue) popMiddlewareFunc(slice []MiddlewareCallback, index int) ([]MiddlewareCallback, MiddlewareCallback) {
	item := slice[index]
	slice = append(slice[:index], slice[index+1:]...)

	return slice, item
}

func (mq *MiddlewareQueue) executeQueuedHandler(w http.ResponseWriter, r *http.Request) {
	var middlewareFunc MiddlewareCallback

	mq.Queue, middlewareFunc = mq.popMiddlewareFunc(mq.Queue, 0)

	handlerFunc := middlewareFunc(mq)
	handlerFunc.ServeHTTP(w, r)
}

func (mq *MiddlewareQueue) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mq.executeQueuedHandler(w, r)
}

func (mq *MiddlewareQueue) Execute() {
	mq.executeQueuedHandler(mq.ResponseWriter, mq.Request)
}
