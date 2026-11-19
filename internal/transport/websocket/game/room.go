package game

import (
	"ca/internal/model/enum"
	"ca/internal/room"
	"ca/internal/transport/dto"
	"ca/internal/transport/websocket"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
)

type HandlerFunc func(ctx context.Context, action PlayerAction)

type Room struct {
	Id      string
	Players []*websocket.Client
	State   *room.State
	Actions chan PlayerAction

	context context.Context

	handlers map[string]HandlerFunc
}

type PlayerAction struct {
	Client  *websocket.Client
	Message websocket.IncomingMessage
}

func NewRoom(
	id string,
	state *room.State) *Room {
	r := &Room{
		Id:      id,
		Players: make([]*websocket.Client, 0, 2),
		Actions: make(chan PlayerAction, 16),
		State:   state,
	}

	r.handlers = map[string]HandlerFunc{
		"move": r.handleMoveAction,

		"dice_roll":        r.handleDiceRollAction,
		"dice_thrown_flag": r.handleDiceThrownFlagAction,

		"tree_chop": r.handleTreeChopAction,

		"end_turn_phase": r.handleEndTurnPhase,

		"draw_card":    r.handleDrawCardAction,
		"card_to_unit": r.handleCardToUnitAction,

		"take_chest":          r.handleTakeSecretAction,
		"take_chest_resource": r.handleTakeSecretResourceAction,
	}

	return r
}

func (r *Room) handleAction(ctx context.Context, action PlayerAction) {
	key := action.Message.Type
	if key == "action" {
		key = action.Message.ActionType
	}
	fmt.Println(string(action.Message.Raw))
	handler, ok := r.handlers[key]

	if !ok {
		r.sendError(action.Client, "unknown action")
		return
	}

	handler(ctx, action)
}

func (r *Room) handleDiceRollAction(ctx context.Context, action PlayerAction) {
	var payload DiceRollPayload
	if err := json.Unmarshal(action.Message.Raw, &payload); err != nil {
		r.sendError(action.Client, "Invalid payload\n"+err.Error())
		return
	}

	playerId := action.Client.PlayerId
	if len(playerId) == 0 {
		r.sendError(action.Client, "no found playerId")
		return
	}

	payload.Id = playerId
	r.sendAllFlat("state_patch", payload, "dice_roll")
}

func RollDices(maxCost int, count int) []int {
	results := make([]int, count)
	for i := range results {
		results[i] = rand.Intn(int(maxCost)) + 1
	}
	return results
}

func (r *Room) handleDiceThrownFlagAction(ctx context.Context, action PlayerAction) {
	delete(r.State.Deck.DrawCounts, action.Client.PlayerId)

	dices := RollDices(6, 5)
	countDuplicates := func(diceValues []int) int8 {
		counts := make(map[int]int)
		for _, v := range diceValues {
			counts[v]++
		}
		duplicates := int8(0)
		for _, c := range counts {
			if c > 1 {
				duplicates += int8(c - 1)
			}
		}
		return duplicates
	}

	duplicates := countDuplicates(dices)
	r.State.Deck.AddDrawCount(action.Client.PlayerId, "", duplicates)

	r.sendAllFlat("state_patch", DiceThrownFlagPayload{Id: action.Client.PlayerId, DiceValues: dices},
		"dice_thrown_flag")
}

func (r *Room) handleMoveAction(ctx context.Context, action PlayerAction) {
	var payload MovePayload
	if err := json.Unmarshal(action.Message.Raw, &payload); err != nil {
		r.sendError(action.Client, "Invalid payload\n"+err.Error())
		return
	}

	playerId := action.Client.PlayerId
	if len(playerId) == 0 {
		r.sendError(action.Client, "no found playerId")
		return
	}

	unit := r.State.Game.Players[playerId].Units[payload.UnitId]

	//TODO сделать проверку, что все куда он наступает является пустой клеткой,кроме последней
	//for _, point := payload.Path{
	//	board := r.State.Game.Board[]
	//}

	lastPoint := payload.Path[len(payload.Path)-1]
	unit.PositionQ = lastPoint.Q
	unit.PositionR = lastPoint.R

	r.sendAllFlat("state_patch", payload, "unit_move")
}

func (r *Room) handleTreeChopAction(ctx context.Context, action PlayerAction) {
	var payload TreeChopPayload
	if err := json.Unmarshal(action.Message.Raw, &payload); err != nil {
		r.sendError(action.Client, "Invalid payload\n"+err.Error())
		return
	}
	playerId := action.Client.PlayerId
	if len(playerId) == 0 {
		r.sendError(action.Client, "no found playerId")
		return
	}

	lastPoint := payload.Path[len(payload.Path)-1]
	key := fmt.Sprintf("%d,%d", lastPoint.Q, lastPoint.R)
	cell := r.State.Game.Board[key]
	cell.Content = enum.ContentTypeNone

	payload.Id = payload.UnitId
	r.sendAllFlat("state_patch", payload, "tree_chop")

	r.State.Deck.AddDrawCount(playerId, payload.UnitId, 1)

	//создаем сам предмет и добавляем в хранилище
	err := r.State.SetupItemByTemplate(int16(enum.ItemTemplateIdWood), playerId, "")
	if err != nil {
		r.sendError(action.Client, "SetupItemByTemplate error")
		return
	}

	itemsDTO := dto.GameItemsToDTO(r.State.Game.Players[playerId].Storage)
	storagePayload := PlayerStoragePayload{
		PlayerId: playerId,
		Items:    itemsDTO,
	}

	r.sendAllFlat("state_patch", storagePayload, "player_storage")
}

func (r *Room) handleEndTurnPhase(ctx context.Context, action PlayerAction) {
	var payload EndTurnPhasePayload
	if err := json.Unmarshal(action.Message.Raw, &payload); err != nil {
		r.sendError(action.Client, "Invalid payload\n"+err.Error())
		return
	}
	current, ok := r.State.PhaseMap[payload.Phase]
	if !ok {
		r.sendError(action.Client, "unknown phase")
		return
	}

	next := current + 1

	switch next {
	case room.TurnPhaseDrawCards:
		playerDraws, ok := r.State.Deck.DrawCounts[action.Client.PlayerId]
		if !ok {
			playerDraws = map[string]int8{}
		}

		resultPayload := &TurnPhasePayload{
			Id:    action.Client.PlayerId,
			Phase: r.State.PhaseMapToString[next],
			Draws: playerDraws,
		}
		r.sendAllFlat("state_patch", resultPayload, "turn_phase")
	case room.TurnPhaseRevealSecrets:
		break
	}
}

func (r *Room) handleDrawCardAction(ctx context.Context, action PlayerAction) {
	var payload DrawCardPayload
	if err := json.Unmarshal(action.Message.Raw, &payload); err != nil {
		r.sendError(action.Client, "Invalid payload\n"+err.Error())
		return
	}

	card, ok := r.State.Deck.DrawCard(action.Client.PlayerId, payload.UnitId)
	if !ok {
		r.sendError(action.Client, "DrawCard error")
	}

	playerDraws, ok := r.State.Deck.DrawCounts[action.Client.PlayerId]
	if !ok {
		playerDraws = map[string]int8{}
	}

	resultPayload := &DrawCardResultPayload{
		Id:     action.Client.PlayerId,
		UnitId: payload.UnitId,
		Card: CardPayload{
			Suit:       card.Suit,
			BaseRank:   card.BaseRank,
			BonusValue: 0,
			UnitId:     payload.UnitId,
		},
		Draws: playerDraws,
	}
	r.sendAllFlat("state_patch", resultPayload, "draw_card")

	if card.Suit == 5 || card.Suit == 6 {
		itemId := int16(enum.ItemTemplateIdRedJoker)
		if card.Suit == 6 {
			itemId = int16(enum.ItemTemplateIdBlackJoker)
		}

		playerId := action.Client.PlayerId
		err := r.State.SetupItemByTemplate(itemId, playerId, payload.UnitId)
		if err != nil {
			r.sendError(action.Client, "SetupItemByTemplate error")
			return
		}

		itemsDTO := dto.GameItemsToDTO(r.State.Game.Players[playerId].Storage)
		storagePayload := PlayerStoragePayload{
			PlayerId: playerId,
			Items:    itemsDTO,
		}

		r.sendAllFlat("state_patch", storagePayload, "player_storage")
	}
}

func (r *Room) handleCardToUnitAction(ctx context.Context, action PlayerAction) {
	var payload CardToUnitPayload
	if err := json.Unmarshal(action.Message.Payload, &payload); err != nil {
		r.sendError(action.Client, "Invalid payload\n"+err.Error())
		return
	}

	player := r.State.Game.Players[action.Client.PlayerId]
	css := room.CardSuit(payload.Suit)
	//
	change := int16(payload.BonusValue)

	if payload.IsPositive {
		change += int16(payload.BaseRank)
	} else {
		change -= int16(payload.BaseRank)
	}
	//

	unit := player.Units[payload.UnitId]
	payloadResult := UnitUpdate{
		PlayerId: action.Client.PlayerId,
		UnitId:   payload.UnitId,
	}
	switch css {
	case room.CardSuitHP:
		unit.MaxHp += change
		unit.Hp += change
		payloadResult.MaxHp = &unit.MaxHp
		payloadResult.Hp = &unit.Hp

	case room.CardSuitDefense:
		unit.Defense += change
		payloadResult.Defense = &unit.Defense

	case room.CardSuitAttack:
		unit.Attack += change
		payloadResult.Attack = &unit.Attack

	case room.CardSuitLuck:
		unit.Luck += change
		payloadResult.Luck = &unit.Luck

	case room.CardSuitSpeed:
		unit.Speed += change
		payloadResult.Speed = &unit.Speed
	default:
		r.sendError(action.Client, "Invalid suit")
	}

	r.sendAllWithEntity("state_patch", "unit_update", payloadResult)
}

func (r *Room) handleTakeSecretAction(ctx context.Context, action PlayerAction) {
	var payload TakeSecretPayload
	if err := json.Unmarshal(action.Message.Raw, &payload); err != nil {
		r.sendError(action.Client, "Invalid payload\n"+err.Error())
		return
	}

	if payload.RequestId == "" {
		r.sendError(action.Client, "no found request id")
		return
	}

	pick := r.State.StandardSecretReward.PickSecret()

	resultPayload := TakeSecretResultPayload{
		Id:        action.Client.PlayerId,
		UnitId:    payload.UnitId,
		RequestId: payload.RequestId,
		Reward:    int16(pick.Resource),
	}

	r.State.StandardSecretReward.Cache[payload.RequestId] = &pick

	r.sendAllFlat("state_patch", resultPayload, "take_chest")
}

func (r *Room) handleTakeSecretResourceAction(ctx context.Context, action PlayerAction) {
	var payload TakeSecretResourcePayload
	if err := json.Unmarshal(action.Message.Raw, &payload); err != nil {
		r.sendError(action.Client, "Invalid payload\n"+err.Error())
		return
	}

	if payload.RequestId == "" {
		r.sendError(action.Client, "no found request id")
		return
	}

	reward := r.State.StandardSecretReward.Cache[payload.RequestId]
	if reward == nil {
		r.sendError(action.Client, "no found reward on cache")
		return
	}

	itemTemplate := r.State.ItemTemplates[int16(reward.Resource)]
	if itemTemplate == nil {
		r.sendError(action.Client, "no found item template")
		return
	}

	if itemTemplate.Id == int16(enum.ItemTemplateIdBuff) ||
		itemTemplate.Id == int16(enum.ItemTemplateIdDebuff) {
		//TODO будет показ колоды, сразу

		return
	}

	playerId := action.Client.PlayerId
	err := r.State.SetupItemByTemplate(itemTemplate.Id, playerId, payload.UnitId)
	if err != nil {
		r.sendError(action.Client, "SetupItemByTemplate error")
		return
	}

	if r.State.IsPlayerItem(itemTemplate.Id) {
		itemsDTO := dto.GameItemsToDTO(r.State.Game.Players[playerId].Storage)
		storagePayload := PlayerStoragePayload{
			PlayerId: playerId,
			Items:    itemsDTO,
		}

		r.sendAllFlat("state_patch", storagePayload, "player_storage")
	} else if r.State.IsUnitItem(itemTemplate.Id) {
		itemsDTO := dto.GameItemsToDTO(r.State.Game.Players[playerId].Units[payload.UnitId].Pocket)
		storagePayload := UnitStoragePayload{
			PlayerId: playerId,
			UnitId:   payload.UnitId,
			Items:    itemsDTO,
		}
		r.sendAllFlat("state_patch", storagePayload, "unit_item")
	}
}

func (r *Room) Run(ctx context.Context) {
	for action := range r.Actions {
		r.handleAction(ctx, action)
	}
}

func (r *Room) AddPlayer(client *websocket.Client) {
	log.Printf("AddPlayer called, current players: %d, userId: %s", len(r.Players), client.UserId)

	if r.IsFull() {
		r.sendError(client, "already full")
		return
	}

	for _, p := range r.Players {
		if p.UserId == client.UserId {
			r.sendAll("state_snapshot", dto.ToGameDTO(r.State.Game, r.State.ItemTemplates))
			return
		}
	}

	r.Players = append(r.Players, client)
	for _, p := range r.State.Game.Players {
		if p.UserId == client.UserId {
			client.PlayerId = p.PlayerId
			r.sendTo(client, "welcome", websocket.WelcomeMessage{
				PlayerId: p.PlayerId,
				Seat:     p.Seat,
			})
		}
	}

	if r.IsFull() {
		r.sendAll("state_snapshot", dto.ToGameDTO(r.State.Game, r.State.ItemTemplates))
		return
	}
}

func (r *Room) RemovePlayer(client *websocket.Client) {
	for i, p := range r.Players {
		if p == client {
			r.Players = append(r.Players[:i], r.Players[i+1:]...)
			close(client.Send)
			break
		}
	}
}

func (r *Room) IsEmpty() bool {
	return len(r.Players) == 0
}

func (r *Room) IsFull() bool { return len(r.Players) == 2 }

func (r *Room) sendError(client *websocket.Client, msg string) {
	data, _ := json.Marshal(websocket.OutgoingMessage{
		Type:    "error",
		Payload: map[string]any{"message": msg},
	})
	client.Send <- data
}

func (r *Room) sendAll(msgType string, payload any) {
	log.Printf("sendAll: %s, players: %d", msgType, len(r.Players))
	msg, _ := json.Marshal(websocket.OutgoingMessage{Type: msgType, Payload: payload})
	for _, p := range r.Players {
		if p != nil {
			p.Send <- msg
		}
	}
}

func (r *Room) sendAllWithEntity(msgType string, entity string, payload any) {
	log.Printf("sendAll: %s, players: %d", msgType, len(r.Players))
	msg, _ := json.Marshal(websocket.OutgoingMessage{Type: msgType, Entity: entity, Payload: payload})
	for _, p := range r.Players {
		if p != nil {
			p.Send <- msg
		}
	}
}

func (r *Room) sendAllFlat(msgType string, payload any, entity string) {
	raw, _ := json.Marshal(payload)
	var flat map[string]any
	json.Unmarshal(raw, &flat)
	flat["type"] = msgType
	flat["entity"] = entity

	msg, _ := json.Marshal(flat)
	for _, p := range r.Players {
		if p != nil {
			p.Send <- msg
		}
	}
	log.Println(string(msg))
}

func (r *Room) sendTo(client *websocket.Client, msgType string, payload any) {
	data, _ := json.Marshal(websocket.OutgoingMessage{Type: msgType, Payload: payload})
	client.Send <- data
}
