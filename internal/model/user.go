package model

import (
	"ca/internal/model/enum"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id uuid.UUID `db:"id"`

	Username string    `db:"username"`
	Password string    `db:"password"`
	Email    string    `db:"email"`
	Role     enum.Role `db:"role"`

	Nickname string `db:"nickname"`

	Wins   int `db:"wins"`
	Draws  int `db:"draws"`
	Losses int `db:"losses"`
	Rating int `db:"rating"`

	TotalBattlesWithPlayers    int `db:"total_battles_with_players"`
	TotalBattlesWithPlayersWon int `db:"total_battles_with_players_won"`

	TotalBattlesWithMonsters    int `db:"total_battles_with_monsters"`
	TotalBattlesWithMonstersWon int `db:"total_battles_with_monsters_won"`

	TotalCharactersLost   int `db:"total_characters_lost"`
	TotalCharactersKilled int `db:"total_characters_killed"`

	TotalTreesChopped int `db:"total_trees_chopped"`
	TotalTreesPlanted int `db:"total_trees_planted"`

	TotalSecretsOpened int `db:"total_secrets_opened"`
	TotalCampsPlaced   int `db:"total_camps_placed"`

	TotalDisastersSurvived int `db:"total_disasters_survived"`
	TotalCriticalHits      int `db:"total_critical_hits"`

	TotalTurnsPlayed int `db:"total_turns_played"`
	TotalGamesPlayed int `db:"total_games_played"`

	Deleted   bool      `db:"deleted"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
