package rest

import (
	"ca/internal/service"
)

type GameHandler struct {
	gameUnitService   *service.GameUnitService
	gamePlayerService *service.GamePlayerService
}

func NewGameHandler(gameUnitService *service.GameUnitService, gamePlayerService *service.GamePlayerService) *GameHandler {
	return &GameHandler{gameUnitService: gameUnitService}
}
