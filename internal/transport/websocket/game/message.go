package game

import (
	"ca/internal/transport/dto"
)

type PathPoint struct {
	Q int16 `json:"q"`
	R int16 `json:"r"`
}

type MovePayload struct {
	RequestId string      `json:"request_id"`
	UnitId    string      `json:"unitId"`
	Path      []PathPoint `json:"path"`
	DiceIndex int         `json:"diceIndex"`
}

type DiceRollPayload struct {
	RequestId string `json:"requestId"`
	Id        string `json:"id"`
}

type DiceThrownFlagPayload struct {
	Id         string `json:"id"`
	DiceValues []int  `json:"diceValues"`
}

type TreeChopPayload struct {
	RequestId string      `json:"requestId"`
	UnitId    string      `json:"unitId"`
	Path      []PathPoint `json:"path"`
	DiceIndex int         `json:"diceIndex"`
	Id        string      `json:"id,omitempty"`
}

type EndTurnPhasePayload struct {
	Phase string `json:"phase"`
}

type TurnPhasePayload struct {
	Id    string          `json:"id"`
	Phase string          `json:"phase"`
	Draws map[string]int8 `json:"draws"`
}

type DrawCardPayload struct {
	UnitId string `json:"unitId"`
}

type CardPayload struct {
	Suit       int    `json:"suit"`
	BaseRank   int    `json:"baseRank"`
	BonusValue int    `json:"bonusValue"`
	UnitId     string `json:"unitId"`
}

type DrawCardResultPayload struct {
	Id     string          `json:"id"` // playerId
	UnitId string          `json:"unitId"`
	Card   CardPayload     `json:"card"`
	Draws  map[string]int8 `json:"draws"`
}

type CardToUnitPayload struct {
	Suit       int    `json:"suit"`
	BaseRank   int    `json:"baseRank"`
	BonusValue int    `json:"bonusValue"`
	UnitId     string `json:"unitId"`
	IsPositive bool   `json:"isPositive"`
}

type TakeSecretPayload struct {
	RequestId string `json:"requestId"`
	UnitId    string `json:"unitId"`
}

type TakeSecretResultPayload struct {
	Id        string `json:"id"`
	UnitId    string `json:"unitId"`
	RequestId string `json:"requestId"`
	Reward    int16  `json:"reward"`
}

type TakeSecretResourcePayload struct {
	RequestId string `json:"requestId"`
	UnitId    string `json:"unitId"`
}

type PlayerStoragePayload struct {
	PlayerId string             `json:"id"`
	Items    []*dto.GameItemDTO `json:"items"`
}

type UnitStoragePayload struct {
	PlayerId string             `json:"id"`
	UnitId   string             `json:"unitId"`
	Items    []*dto.GameItemDTO `json:"items"`
}

type UnitUpdate struct {
	PlayerId string `json:"player_id"`
	UnitId   string `json:"unit_id"`

	Hp      *int16 `json:"hp,omitempty"`
	MaxHp   *int16 `json:"max_hp,omitempty"`
	Attack  *int16 `json:"attack,omitempty"`
	Defense *int16 `json:"defense,omitempty"`
	Luck    *int16 `json:"luck,omitempty"`
	Speed   *int16 `json:"speed,omitempty"`
}
