package main

import (
	"ca/internal/config"
	"ca/internal/repository/postgres"
	"ca/internal/service"
	"ca/internal/transport/rest"
	"ca/pkg/jwtutil"
	"context"
	"log"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load("config/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.Database.Postgres.Url)
	if err != nil {
		log.Fatal(err.Error())
	}

	authRepository := postgres.NewAuthRepository(pool)
	userRepository := postgres.NewUserRepository(pool)
	gameRepository := postgres.NewGameRepository(pool)
	gamePlayerRepository := postgres.NewGamePlayerRepository(pool)
	gameUnitRepository := postgres.NewGameUnitRepository(pool)
	gameCellRepository := postgres.NewGameCellRepository(pool)
	gameItemRepository := postgres.NewGameItemRepository(pool)
	itemTemplateRepository := postgres.NewItemTemplateRepository(pool)

	jwtUtil := jwtutil.NewJwtUtil(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL)
	authService := service.NewAuthService(authRepository, jwtUtil)
	gameCellService := service.NewGameCellService(gameCellRepository)
	gameService := service.NewGameService(gameRepository, gameCellService)
	gamePlayerService := service.NewGamePlayerService(gamePlayerRepository)
	gameUnitService := service.NewGameUnitService(gameUnitRepository)
	gameItemService := service.NewGameItemService(gameItemRepository)
	userService := service.NewUserService(userRepository)
	matchmakingService := service.NewMatchmakingService(gameService)
	gameSessionService := service.NewGameSessionService(gameRepository, gameCellRepository, gamePlayerRepository,
		gameUnitRepository, itemTemplateRepository)

	server := rest.New(
		gin.Default(),
		authService,
		gameService,
		gamePlayerService,
		gameUnitService,
		userService,
		matchmakingService,
		gameSessionService,
		gameItemService,
		slog.Default(),
		jwtUtil)
	server.Handle()

	log.Fatal(server.ListenAndServe(cfg.Host))
}
