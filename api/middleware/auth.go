package middleware

import (
	"go-rest-api/api/response"
	"go-rest-api/utils"
	"net/http"
)

func AuthenticatedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtToken := r.Header.Get(utils.AuthorizationKey)
		if jwtToken != "" {
			user, err := utils.GetUserByJwtToken(jwtToken)
			if err != nil {
				_ = response.ServeJSON(w, http.StatusUnauthorized, nil, nil, "user could not be authenticated", nil)
				return
			}
			r.Header.Set(utils.AuthorizationKey, user.Username)
		}
	})
}
