-- 000001_games.up.sql

CREATE TABLE games
(
    id                     UUID PRIMARY KEY     DEFAULT gen_random_uuid(),

    owner_id               UUID                 DEFAULT NULL,

    status                 SMALLINT    NOT NULL DEFAULT 0,

    winner_id              UUID,
    rating_delta           SMALLINT    NOT NULL DEFAULT 0,

    board_radius           SMALLINT    NOT NULL DEFAULT 0,
    seeds                  TEXT        NOT NULL DEFAULT '',

    turn_count             INT         NOT NULL DEFAULT 0,
    current_turn_player_id VARCHAR(16) NOT NULL DEFAULT '',
    number_characters      INT         NOT NULL DEFAULT 5,
    number_cubes           INT         NOT NULL DEFAULT 5,

    current_disaster       SMALLINT    NOT NULL DEFAULT 0,
    disaster_started_at    INT         NOT NULL DEFAULT 0,
    disaster_zone          SMALLINT    NOT NULL DEFAULT 0,

    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at            TIMESTAMPTZ
);