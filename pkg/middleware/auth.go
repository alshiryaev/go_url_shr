package middleware

import (
	"fmt"
	"go_purple/configs"
	"go_purple/pkg/jwt"
	"net/http"
	"strings"
)

func IsAuth(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		token := strings.TrimPrefix(auth, "Bearer ")
		isValid, data := jwt.NewJWT(config.Auth.Secret).Parse(token)
		fmt.Println("isValid=", isValid)
		fmt.Println("data=", data)
		next.ServeHTTP(w, r)
	})
}
