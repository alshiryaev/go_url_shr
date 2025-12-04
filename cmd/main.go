package main

import (
	"fmt"
	"go_purple/configs"
	"go_purple/internal/auth"
	"go_purple/internal/link"
	"go_purple/internal/stat"
	"go_purple/internal/user"
	"go_purple/pkg/db"
	"go_purple/pkg/middleware"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()

	db := db.NewDb(conf)
	router := http.NewServeMux()

	// Repositories
	linkReposotory := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)
	statRepositoty := stat.NewStatRepository(db)

	// Services
	authService := auth.NewAuthService(userRepository)

	// Handler
	auth.NewAuthHandler(router, auth.AuthHanderDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkReposotory,
		StatRepositoty: statRepositoty,
		Config:         conf,
	})

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
