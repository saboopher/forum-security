package handler

import (
	"context"
	"net/http"
)

func SetErrorObject(r *http.Request, errObject interface{}) *http.Request {
	ctx := context.WithValue(r.Context(), "error", errObject)
	return r.WithContext(ctx)
}
