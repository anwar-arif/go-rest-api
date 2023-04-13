package middleware

import (
	"go-rest-api/api/response"
	"net/http"
)

func AppKeyChecker(appKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("Application-Key")
			if len(key) == 0 {
				response.Serve(w, http.StatusUnauthorized, "'Application-Key' required", "")
				return
			}
			if key != appKey {
				response.Serve(w, http.StatusUnauthorized, "invalid 'Application-Key'", "")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
