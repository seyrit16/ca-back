-- 000005_game_units.up.sql

CREATE TABLE game_units
(
    game_id   UUID        NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    player_id VARCHAR(16) NOT NULL,
    unit_id   VARCHAR(16) NOT NULL,

    seat       SMALLINT NOT NULL DEFAULT 0,
    position_q SMALLINT NOT NULL DEFAULT 0,
    position_r SMALLINT NOT NULL DEFAULT 0,

    character_template_id SMALLINT NOT NULL DEFAULT 0,

    name VARCHAR(64) NOT NULL DEFAULT '',

    hp      SMALLINT NOT NULL DEFAULT 0,
    max_hp  SMALLINT NOT NULL DEFAULT 0,
    attack  SMALLINT NOT NULL DEFAULT 0,
    defense SMALLINT NOT NULL DEFAULT 0,
    luck    SMALLINT NOT NULL DEFAULT 0,
    speed   SMALLINT NOT NULL DEFAULT 0,

    critical_chance SMALLINT NOT NULL DEFAULT 0,
    critical_damage SMALLINT NOT NULL DEFAULT 0,

    accuracy SMALLINT NOT NULL DEFAULT 0,
    dodge    SMALLINT NOT NULL DEFAULT 0,

    equip_head   SMALLINT,
    equip_body   SMALLINT,
    equip_boots  SMALLINT,
    equip_weapon SMALLINT,

    is_alive BOOLEAN NOT NULL DEFAULT TRUE,

    PRIMARY KEY (game_id, player_id, unit_id)
);