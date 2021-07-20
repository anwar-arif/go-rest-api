package middleware

import (
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go-rest-api/api/response"
	"go-rest-api/utils"
	"log"
	"net/http"
)

// ContextKey hold the key of a context
type ContextKey string

// List of contexts
const (
	UserContext ContextKey = "user"
)

/*
func GetUser(r *http.Request) *service.UserData {
	v := r.Context().Value(UserContext)
	if v == nil {
		panic(errors.New("middleware: GetUser called without calling auth middleware prior"))
	}
	u, _ := v.(*service.UserData)
	return u
}

// Auth returns authentication middleware
func Auth(auth *service.Auth) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tok := r.Header.Get("Authorization")
			if !strings.HasPrefix(tok, "Bearer ") {
				resp.ServeUnauthorized(w, r, errors.New("unauthorized"))
				return
			}
			tok = strings.TrimSpace(strings.TrimPrefix(tok, "Bearer "))
			if tok == "" {
				resp.ServeUnauthorized(w, r, errors.New("unauthorized"))
				return
			}
			_, u, err := auth.Check(tok)
			if err != nil {
				if err == service.ErrUserNotFound ||
					err == service.ErrUserDisabled {
					resp.ServeUnauthorized(w, r, errors.New("unauthorized"))
					return
				}
				resp.ServeInternalServerError(w, r, err)
				return
			}
			if u == nil {
				resp.ServeUnauthorized(w, r, errors.New("unauthorized"))
				return
			}

			ctx := context.WithValue(r.Context(), UserContext, u)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}*/

//func AuthenticatedOnly(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		jwtTkn := r.Header.Get(utils.AuthorizationKey)
//		if jwtTkn == "" {
//			_ = response.ServeJSON(w, http.StatusUnauthorized, nil, nil, "Missing authorization token", nil)
//			return
//		}
//
//		user, err := utils.GetUserInfoByTokenFromAuth(jwtTkn)
//		if err != nil {
//			_ = response.ServeJSON(w, http.StatusUnauthorized, nil, nil, "User could not be authenticated", nil)
//			return
//		}
//
//		r.Header.Set(utils.AdminUserKey, user.Username)
//		r.Header.Set(utils.RoleKey, user.Role)
//
//		requestId := r.Header.Get(middleware.RequestIDHeader)
//		if requestId == "" {
//			requestId = uuid.New().String()
//		}
//
//		ip := r.Header.Get(utils.RealUserIpKey)
//		if ip == "" {
//			ip = "0.0.0.0"
//		}
//
//		next.ServeHTTP(w, r)
//	})
//}

// AuthenticatedOnly Authenticates by token or secret key
func AuthenticatedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Token authentication
		jwtTkn := r.Header.Get(utils.AuthorizationKey)
		if jwtTkn != "" {
			user, err := utils.GetUserInfoByTokenFromAuth(jwtTkn)
			if err != nil {
				_ = response.ServeJSON(w, http.StatusUnauthorized, nil, nil, "User could not be authenticated", nil)
				return
			}

			r.Header.Set(utils.AdminUserKey, user.Username)
			r.Header.Set(utils.RoleKey, user.Role)

			requestId := r.Header.Get(middleware.RequestIDHeader)
			if requestId == "" {
				requestId = uuid.New().String()
			}

			ip := r.Header.Get(utils.RealUserIpKey)
			if ip == "" {
				ip = "0.0.0.0"
			}

			next.ServeHTTP(w, r)
		} else {
			// Secret key authentication
			secretKey := r.Header.Get(utils.KeyForSecretKey)
			if secretKey == "" {
				log.Println("Missing authorization secret")
				_ = response.ServeJSON(w, http.StatusUnauthorized, nil, nil, "Missing authorization token or secret", nil)
				return
			}

			viper.AutomaticEnv()
			actualKey := viper.GetString("app.api_secret_key")
			if actualKey != secretKey {
				_ = response.ServeJSON(w, http.StatusUnauthorized, nil, nil, "Incorrect authorization secret", nil)
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
		}
	})
}

//func SecretOnly(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		secretKey := r.Header.Get(utils.KeyForSecretKey)
//		if secretKey == "" {
//			log.Println("Missing authorization secret")
//			_ = response.ServeJSON(w, http.StatusUnauthorized, nil, nil, "Missing authorization secret", nil)
//			return
//		}
//
//		viper.AutomaticEnv()
//		actualKey := viper.GetString("app.api_secret_key")
//		if actualKey != secretKey {
//			_ = response.ServeJSON(w, http.StatusUnauthorized, nil, nil, "Incorrect authorization secret", nil)
//			return
//		}
//
//		requestId := r.Header.Get(middleware.RequestIDHeader)
//		if requestId == "" {
//			requestId = uuid.New().String()
//			r.Header.Set(middleware.RequestIDHeader, requestId)
//		}
//
//		ip := r.Header.Get(utils.RealUserIpKey)
//		if ip == "" {
//			ip = "0.0.0.0"
//			r.Header.Set(utils.RealUserIpKey, ip)
//		}
//
//		next.ServeHTTP(w, r)
//	})
//}
