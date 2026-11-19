-- 000003_users.up.sql

CREATE TABLE users
(
    id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    username VARCHAR(64)  NOT NULL UNIQUE,
    password VARCHAR(256) NOT NULL,
    email    VARCHAR(256) NOT NULL UNIQUE,
    role     VARCHAR(32)  NOT NULL DEFAULT 'user',

    nickname VARCHAR(64)  NOT NULL DEFAULT '',

    wins   INT NOT NULL DEFAULT 0,
    draws  INT NOT NULL DEFAULT 0,
    losses INT NOT NULL DEFAULT 0,
    rating INT NOT NULL DEFAULT 0,

    total_battles_with_players     INT NOT NULL DEFAULT 0,
    total_battles_with_players_won INT NOT NULL DEFAULT 0,

    total_battles_with_monsters     INT NOT NULL DEFAULT 0,
    total_battles_with_monsters_won INT NOT NULL DEFAULT 0,

    total_characters_lost   INT NOT NULL DEFAULT 0,
    total_characters_killed INT NOT NULL DEFAULT 0,

    total_trees_chopped INT NOT NULL DEFAULT 0,
    total_trees_planted INT NOT NULL DEFAULT 0,

    total_secrets_opened INT NOT NULL DEFAULT 0,
    total_camps_placed   INT NOT NULL DEFAULT 0,

    total_disasters_survived INT NOT NULL DEFAULT 0,
    total_critical_hits      INT NOT NULL DEFAULT 0,

    total_turns_played INT NOT NULL DEFAULT 0,
    total_games_played INT NOT NULL DEFAULT 0,

    deleted    BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);