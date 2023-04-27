package middleware

import (
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go-rest-api/api/response"
	"go-rest-api/utils"
	"log"
	"net/http"
	"strings"
)

func AuthenticatedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get(utils.AuthorizationKey)
		secretKey := r.Header.Get(utils.KeyForSecretKey)

		// token authentication
		if bearerToken != "" {
			if !strings.HasPrefix(bearerToken, "Bearer") {
				log.Println("authorization token is missing")
				_ = response.Serve(w, http.StatusUnauthorized, "invalid token", nil)
				return
			}

			// TODO: JWTClaimsFromContext throws err
			//jwtToken := strings.TrimPrefix(bearerToken, "Bearer ")
			//apiSecretKey := viper.GetString("app.api_secret_key")
			//
			//claims, err := utils.GetClaimsFromToken(jwtToken, apiSecretKey)
			//if err != nil {
			//	log.Println("error while getClaimsFromToken: ", err.Error())
			//	_ = response.Serve(w, http.StatusUnauthorized, "user could not be authenticated", nil)
			//	return
			//}
			//r.Header.Set(utils.AuthorizationKey, user.Email)
			//r.Header.Set(utils.RoleKey, user.Role)

			requestId := r.Header.Get(middleware.RequestIDHeader)
			if requestId == "" {
				requestId = uuid.New().String()
				r.Header.Set(middleware.RequestIDHeader, requestId)
			}

			ip := r.Header.Get(utils.RealUserIpKey)
			if ip == "" {
				ip = "0.0.0.0"
				r.Header.Set(utils.RealUserIpKey, ip)
			}

			//r = r.WithContext(utils.SetJWTClaimsContext(r.Context(), requestId, *claims))

			next.ServeHTTP(w, r)
		} else if secretKey != "" {
			// secret key authentication

			if secretKey == "" {
				log.Println("missing secret key")
				_ = response.Serve(w, http.StatusUnauthorized, "missing authorization token or secret", nil)
				return
			}

			viper.AutomaticEnv()
			actualKey := viper.GetString("app.api_secret_key")
			if actualKey != secretKey {
				_ = response.Serve(w, http.StatusUnauthorized, "incorrect authorization secret", nil)
				return
			}

			requestId := r.Header.Get(middleware.RequestIDHeader)
			if requestId == "" {
				requestId = uuid.New().String()
				r.Header.Set(middleware.RequestIDHeader, requestId)
			}

			ip := r.Header.Get(utils.RealUserIpKey)
			if ip == "" {
				ip = "0.0.0.0"
				r.Header.Set(utils.RealUserIpKey, ip)
			}

			next.ServeHTTP(w, r)
		} else {
			_ = response.Serve(w, http.StatusUnauthorized, "missing authorization token or secret", nil)
			return
		}
	})
}
