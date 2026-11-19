package rest

import (
	"ca/internal/mapgen"
	"ca/internal/service"
	. "ca/internal/transport/websocket/draft"
	"ca/internal/transport/websocket/game"
	"ca/pkg/jwtutil"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type api struct {
	router *gin.Engine

	authService        *service.AuthService
	gameService        *service.GameService
	userService        *service.UserService
	gamePlayerService  *service.GamePlayerService
	gameUnitService    *service.GameUnitService
	matchmakingService *service.MatchmakingService
	gameSessionService *service.GameSessionService
	gameItemService    *service.GameItemService

	gameHub  *game.GameHub
	draftHub *DraftHub

	logger  *slog.Logger
	jwtUtil *jwtutil.JWTUtil
}

func New(
	router *gin.Engine,
	authService *service.AuthService,
	gameService *service.GameService,
	gamePlayerService *service.GamePlayerService,
	gameUnitService *service.GameUnitService,
	userService *service.UserService,
	matchmakingService *service.MatchmakingService,
	gameSessionService *service.GameSessionService,
	gameItemService *service.GameItemService,
	logger *slog.Logger,
	jwtUtil *jwtutil.JWTUtil) *api {
	return &api{
		router:             router,
		authService:        authService,
		gameService:        gameService,
		gamePlayerService:  gamePlayerService,
		gameUnitService:    gameUnitService,
		userService:        userService,
		matchmakingService: matchmakingService,
		gameSessionService: gameSessionService,
		gameItemService:    gameItemService,
		logger:             logger,
		jwtUtil:            jwtUtil,
		draftHub:           NewDraftHub(),
		gameHub:            game.NewGameHub(gameSessionService),
	}
}

func (api *api) Handle() {
	api.router.Use(api.middleware())
	authHandler := NewAuthHandler(api.authService, api.jwtUtil.RefreshTTL)

	testRouter := api.router.Group("/test")
	{
		testRouter.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, "pong")
		})
		testRouter.GET("/map/generate", func(c *gin.Context) {
			generator := mapgen.NewMapGenerator(20, 3)
			generator.Generate()
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(http.StatusOK, generator.ToHTML())
		})

		testRouter.GET("/game/:game_id", func(c *gin.Context) {
			log.Printf("[DEBUG] Hit /test/game/:game_id")
			gameId, err := uuid.Parse(c.Param("game_id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game_id"})
				return
			}

			game, err := api.gameSessionService.GetGameSession(c.Request.Context(), gameId)
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}

			c.Header("Content-Type", "application/json; charset=utf-8")
			c.JSON(http.StatusOK, game)
		})
	}

	api.router.GET("/auth", api.serveAuthPage)

	authRouters := api.router.Group("/api/auth")
	{
		authRouters.POST("/register", api.register)
		authRouters.POST("/login", authHandler.Login)
		authRouters.POST("/refresh", authHandler.Refresh)
		authRouters.POST("/logout", authHandler.Logout)
	}

	matchmakingRouters := api.router.Group("/api", api.JWTMiddleware(api.jwtUtil))
	{
		matchmakingHandler := NewMatchmakingHandler(api.matchmakingService)
		matchmaking := matchmakingRouters.Group("/matchmaking")
		{
			matchmaking.POST("/join", matchmakingHandler.Join)
			matchmaking.DELETE("/leave", matchmakingHandler.Leave)
			matchmaking.GET("/status", matchmakingHandler.Status)
		}
	}

	// Web socket
	wsRouters := api.router.Group("/ws", api.JWTMiddleware(api.jwtUtil))
	{
		draftHandler := NewDraftWSHandler(api.draftHub, api.gameUnitService, api.gamePlayerService)
		wsRouters.GET("/draft/:game_id", draftHandler.Connect)

		gameHandler := game.NewGameWSHandler(api.gameHub)
		wsRouters.GET("/game/:game_id", gameHandler.Connect)
	}
}

func (api *api) serveAuthPage(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(c.Writer, c.Request, "internal/web/auth.html")
}

func (api *api) ListenAndServe(addr string) error {
	go api.gameHub.Run(context.Background())
	return api.router.Run(addr)
}
