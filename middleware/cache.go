package middleware

import "net/http"

// WithCacheControl wraps a handler to set a Cache-Control header
func WithCacheControl(h http.Handler, cacheControl string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cacheControl != "" {
			w.Header().Set("Cache-Control", cacheControl)
		}
		h.ServeHTTP(w, r)
	})
}