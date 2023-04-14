package middleware

//func AuthenticatedOnly(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		jwtToken := r.Header.Get(utils.AuthorizationKey)
//		if jwtToken != "" {
//			user, err := utils.GetUserByJwtToken(jwtToken)
//			if err != nil {
//				_ = response.Serve(w, http.StatusUnauthorized, "user could not be authenticated", nil)
//				return
//			}
//			r.Header.Set(utils.AuthorizationKey, user.UserName)
//			r.Header.Set(utils.RoleKey, user.Role)
//		}
//	})
//}
