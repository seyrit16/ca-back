package room

import (
	"ca/internal/model"
	"ca/internal/model/enum"
)

type State struct {
	Game *model.Game

	ItemTemplates map[int16]*model.ItemTemplate
	itemsIds      int32

	// Temp

	PhaseMap         map[string]TurnPhase
	PhaseMapToString map[TurnPhase]string

	Deck *DeckState

	// secret
	StandardSecretReward *SecretReward
}

type TurnPhase int

const (
	TurnPhaseCharacterMoves TurnPhase = iota
	TurnPhaseDrawCards
	TurnPhaseRevealSecrets
)

func NewState(game *model.Game, itemTemplates map[int16]*model.ItemTemplate) *State {
	return &State{
		Game:          game,
		ItemTemplates: itemTemplates,
		itemsIds:      0,
		PhaseMap: map[string]TurnPhase{
			"character_moves": TurnPhaseCharacterMoves,
			"draw_cards":      TurnPhaseDrawCards,
			"reveal_secrets":  TurnPhaseRevealSecrets,
		},
		PhaseMapToString: map[TurnPhase]string{
			TurnPhaseCharacterMoves: "character_moves",
			TurnPhaseDrawCards:      "draw_cards",
			TurnPhaseRevealSecrets:  "reveal_secrets",
		},
		Deck:                 NewDeckState(),
		StandardSecretReward: buildStandardSecretReward(),
	}
}

func (s *State) GetItemId() int32 {
	id := s.itemsIds
	s.itemsIds++
	return id
}

// Создает предмет по шаблону, сохраняет в хранилище к игроку или юниту
// TODO будет запускать сохранение в базу
func (s *State) SetupItemByTemplate(templateId int16, playerId string, unitId string) error {
	item := &model.GameItem{
		Id:         s.GetItemId(),
		GameId:     s.Game.Id,
		TemplateId: templateId,
		PlayerId:   playerId,
		Location:   s.ItemTemplates[templateId].DefaultLocation,
	}

	player := s.Game.Players[playerId]
	if len(unitId) != 0 && s.ItemTemplates[templateId].DefaultLocation == enum.ItemLocationPocket {
		unit := player.Units[unitId]
		item.UnitId = unitId
		unit.Pocket = append(unit.Pocket, item)
		return nil
	}

	player.Storage = append(player.Storage, item)
	return nil
}

func (s *State) IsPlayerItem(templateId int16) bool {
	return s.ItemTemplates[templateId].DefaultLocation == enum.ItemLocationWarehouse
}

func (s *State) IsUnitItem(templateId int16) bool {
	return s.ItemTemplates[templateId].DefaultLocation == enum.ItemLocationPocket
}
