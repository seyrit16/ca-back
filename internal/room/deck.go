package room

import (
	"crypto/rand"
	"math/big"
)

type CardSuit int8

const (
	CardSuitHP CardSuit = iota
	CardSuitDefense
	CardSuitAttack
	CardSuitLuck
	CardSuitSpeed
	CardSuitRedJoker
	CardSuitBlackJoker
)

type CardState struct {
	Suit     int
	BaseRank int
}

type DeckState struct {
	Deck                 []CardState
	DrawnHistoryAll      []CardState
	DrawnHistoryByPlayer map[string][]CardState
	AvailableIndexes     []int

	// Сколько карт игрок откроет (по id персонажа, null если карта не назначена(игрок сам выбирает))
	// (game_player_id, (game_unit_id, draws_count))
	DrawCounts map[string]map[string]int8
}

func buildDeck() []CardState {
	var deck []CardState
	suits := []int{0, 1, 2, 3} // стандартные черви, крести, вини, буби
	values := []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 15}

	for _, suit := range suits {
		for _, value := range values {
			deck = append(deck, CardState{Suit: suit, BaseRank: value})
		}
	}

	for i := 0; i < 13; i++ {
		deck = append(deck, CardState{Suit: 4, BaseRank: 1}) // скорость
	}

	deck = append(deck, CardState{Suit: 5, BaseRank: -1}) // красный джокер
	deck = append(deck, CardState{Suit: 6, BaseRank: -1}) // черный джокер

	return deck
}

func newIndexes(deck []CardState) []int {
	indexes := make([]int, len(deck))
	for i := range indexes {
		indexes[i] = i
	}
	return indexes
}

func NewDeckState() *DeckState {
	deck := buildDeck()
	return &DeckState{
		Deck:             buildDeck(),
		AvailableIndexes: newIndexes(deck),
		DrawCounts:       make(map[string]map[string]int8),
	}
}

func (d *DeckState) initDeck() {
	d.Deck = buildDeck()
	d.DrawnHistoryAll = nil
	d.AvailableIndexes = newIndexes(d.Deck)
}

func (d *DeckState) AddDrawCount(playerId string, unitId string, count int8) {
	if d.DrawCounts[playerId] == nil {
		d.DrawCounts[playerId] = make(map[string]int8)
	}
	d.DrawCounts[playerId][unitId] += count
}

func (d *DeckState) drawFromDeck() (CardState, bool) {
	if len(d.AvailableIndexes) == 0 {
		d.initDeck()
	}

	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(d.AvailableIndexes))))
	pos := int(n.Int64())

	cardIndex := d.AvailableIndexes[pos]
	d.AvailableIndexes = append(d.AvailableIndexes[:pos], d.AvailableIndexes[pos+1:]...)

	card := d.Deck[cardIndex]
	d.DrawnHistoryAll = append(d.DrawnHistoryAll, card)
	return card, true
}

func (d *DeckState) DrawCard(playerId string, unitId string) (CardState, bool) {
	playerDraws, ok := d.DrawCounts[playerId]
	if !ok {
		return CardState{}, false
	}

	current := playerDraws[unitId]
	if current <= 0 {
		return CardState{}, false
	}

	card, ok := d.drawFromDeck()
	if !ok {
		return CardState{}, false
	}

	if current-1 <= 0 {
		delete(playerDraws, unitId)
	} else {
		playerDraws[unitId] = current - 1
	}

	return card, true
}
