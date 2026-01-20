package http

import "net/http"

func getPathParam(r *http.Request, key string) (string, bool) {
	params, ok := r.Context().Value(paramsContextKey).(map[string]string)
	if !ok {
		return "", false
	}
	value, ok := params[key]
	return value, ok
}
