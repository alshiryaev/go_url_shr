package auth

import (
	"go_purple/configs"
	"go_purple/pkg/jwt"
	"go_purple/pkg/req"
	"go_purple/pkg/response"
	"net/http"
)

type AuthHandler struct {
	*configs.Config
	*AuthService
}

type AuthHanderDeps struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHanderDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			// Здесь мы ничего не делаем, так как вся обратка ошибок происходит в HandleBody
			return
		}
		email, err := handler.AuthService.Login(body.Email, body.Password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		token, err := jwt.NewJWT(handler.Auth.Secret).Create(email)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := LoginResponse{
			Token: token,
		}
		response.Json(w, data, http.StatusOK)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		email, err := handler.AuthService.Register(body.Email, body.Password, body.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		token, err := jwt.NewJWT(handler.Auth.Secret).Create(email)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := RegisterResponse{
			Token: token,
		}
		response.Json(w, data, http.StatusOK)
	}
}
