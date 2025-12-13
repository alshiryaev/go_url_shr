package main

import (
	"fmt"
	"go_purple/configs"
	"go_purple/internal/auth"
	"go_purple/internal/link"
	"go_purple/internal/stat"
	"go_purple/internal/user"
	"go_purple/pkg/db"
	"go_purple/pkg/event"
	"go_purple/pkg/middleware"
	"net/http"
)

func App() http.Handler {
	conf := configs.LoadConfig()

	db := db.NewDb(conf)
	router := http.NewServeMux()

	// Repositories
	linkReposotory := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)
	statRepositoty := stat.NewStatRepository(db)

	eventBus := event.NewEventBus()

	// Services
	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(stat.StatServiceDeps{
		StatRepository: statRepositoty,
		EventBus:       eventBus,
	})

	// Handler
	auth.NewAuthHandler(router, auth.AuthHanderDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkReposotory,
		EventBus:       eventBus,
		Config:         conf,
	})
	stat.NewStatHandler(router, stat.StatHanderDeps{
		StatRepository: statRepositoty,
		Config:         conf,
	})

	go statService.AddClick()

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	return stack(router)
}

func main() {

	app := App()
	server := http.Server{
		Addr:    ":8081",
		Handler: app,
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
