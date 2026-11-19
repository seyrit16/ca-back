-- 000006_item_templates.up.sql

CREATE TABLE item_templates
(
    id               SMALLINT PRIMARY KEY,
    name             TEXT      NOT NULL,
    description      TEXT      NOT NULL,
    default_location SMALLINT  NOT NULL,

    created_at       TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO item_templates (id, name, description, default_location)
VALUES (0::smallint, 'Tree', 'A piece of wood', 0::smallint),
       (1::smallint, 'Buff', '', 0::smallint),
       (2::smallint, 'Debuff', '', 0::smallint),
       (3::smallint, 'Heal', 'Вылечите одного своего персонажа', 0::smallint),
       (4::smallint, 'Buff/debuff', '', 0::smallint),
       (5::smallint, 'Provocation', '', 0::smallint),
       (6::smallint, 'Ego strike', '', 0::smallint),
       (7::smallint, 'Teleport', '', 1::smallint),
       (8::smallint, 'Camp', '', 1::smallint),
       (9::smallint, 'Return_stone', '', 1::smallint),
       (10::smallint, 'Red_joker', '', 0::smallint),
       (11::smallint, 'Black_joker', '', 0::smallint);