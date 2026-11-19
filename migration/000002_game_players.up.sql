-- 000002_game_players.up.sql

CREATE TABLE game_players
(
    game_id   UUID        NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    player_id VARCHAR(16) NOT NULL,
    user_id   UUID        NOT NULL,

    seat      SMALLINT    NOT NULL DEFAULT 0,
    is_winner BOOLEAN     NOT NULL DEFAULT FALSE,

    PRIMARY KEY (game_id, player_id)
);