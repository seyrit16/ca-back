package service

import (
	"ca/internal/mapgen"
	"ca/internal/model"
	"ca/internal/model/enum"
	"ca/internal/repository/postgres"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type GameCellService struct {
	gameCellRepository postgres.GameCellRepository
}

func NewGameCellService(gameCellRepository postgres.GameCellRepository) *GameCellService {
	return &GameCellService{gameCellRepository: gameCellRepository}
}

func (s *GameCellService) CreateAllGameCellsForGame(ctx context.Context, gameId uuid.UUID,
	generator *mapgen.MapGenerator) error {
	if len(generator.Tiles) == 0 {
		return errors.New("generator tiles should not be empty")
	}

	cells := make(map[[2]int16]*model.GameCell)

	// Забиваем все из Items
	for gi, item := range generator.Items {
		if item.Type != enum.BiomeTypeNone {
			for ci, coord := range item.Coordinates {
				cells[coord.Key()] = &model.GameCell{
					GameId:          gameId,
					Q:               coord.Q,
					R:               coord.R,
					GenerationGroup: new(int16(gi)),
					IndexInGroup:    new(int16(ci)),
					Biome:           item.Type,
					Content:         enum.ContentTypeNone,
				}
			}
		} else {
			if item.Content == enum.ContentTypeMonster {
				fmt.Println("Монстры были")
				for ci, coord := range item.Coordinates {
					cells[coord.Key()] = &model.GameCell{
						GameId:          gameId,
						Q:               coord.Q,
						R:               coord.R,
						GenerationGroup: new(int16(gi)),
						IndexInGroup:    new(int16(ci)),
						Biome:           item.Type,
						Content:         generator.Tiles[[2]int16{coord.Q, coord.R}].Content,
						RadiusContent:   &generator.Tiles[[2]int16{coord.Q, coord.R}].RadiusContent,
					}
				}
			}
		}
	}

	//Забиваем остальным из Tiles
	for key, tile := range generator.Tiles {
		if cells[key] != nil { // Пропускаем если уже существует клетка
			continue
		}

		cells[key] = &model.GameCell{
			GameId:        gameId,
			Q:             key[0],
			R:             key[1],
			Biome:         tile.Biome,
			Content:       tile.Content,
			RadiusContent: &tile.RadiusContent,
		}
	}

	if err := s.gameCellRepository.CreateAllGameCellsForGame(ctx, gameId, cells); err != nil {
		return err
	}
	return nil
}
