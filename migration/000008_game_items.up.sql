-- 000007_game_items.up.sql

CREATE TABLE game_items
(
    id          SERIAL PRIMARY KEY,
    game_id     UUID     NOT NULL REFERENCES games (id),
    template_id SMALLINT NOT NULL REFERENCES item_templates (id),
    player_id   TEXT     NOT NULL,
    unit_id     SMALLINT,
    location    SMALLINT NOT NULL
);