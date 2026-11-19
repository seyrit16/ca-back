-- 000006_game_units.up.sql

CREATE TABLE game_cells
(
    game_id          UUID        NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    q                SMALLINT    NOT NULL,
    r                SMALLINT    NOT NULL,

    generation_group SMALLINT,
    index_in_group   SMALLINT,

    biome            SMALLINT NOT NULL,
    cell             SMALLINT NOT NULL,
    radius_content   SMALLINT,

    PRIMARY KEY (game_id, q, r)
);