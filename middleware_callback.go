package go_cake

import "net/http"

type MiddlewareCallback func(http.Handler) http.Handler
