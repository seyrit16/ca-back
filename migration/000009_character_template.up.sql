CREATE TABLE character_templates
(
    id              SMALLINT PRIMARY KEY,
    name            TEXT     NOT NULL,
    hp              SMALLINT NOT NULL,
    max_hp          SMALLINT NOT NULL,
    attack          SMALLINT NOT NULL,
    defense         SMALLINT NOT NULL,
    luck            SMALLINT NOT NULL,
    speed           SMALLINT NOT NULL,
    critical_chance SMALLINT NOT NULL,
    critical_damage SMALLINT NOT NULL,
    accuracy        SMALLINT NOT NULL,
    dodge           SMALLINT NOT NULL
);

INSERT INTO character_templates(id, name, hp, max_hp, attack, defense, luck, speed, critical_chance, critical_damage, accuracy, dodge)
VALUES (0, 'Аргул', 1,1, 5,0,0,1,50,50,50,50),
       (1, 'Асока', 3,3, 3,1,0,2,50,50,50,50),
       (2, 'Контроллер', 5,5, 3,5,0,1,50,50,50,50),
       (3, 'Мегатрон', 5,5, 3,5,0,1,50,50,50,50),
       (4, 'Нуб Сайбот', 5,5, 3,5,0,1,50,50,50,50),
       (5, 'Сакура', 5,5, 3,5,0,1,50,50,50,50),
       (6, 'Стенд', 5,5, 3,5,0,1,50,50,50,50),
       (7, 'Хантер', 5,5, 3,5,0,1,50,50,50,50),
       (8, 'Цирик', 5,5, 3,5,0,1,50,50,50,50),
       (9, 'Эцио', 5,5, 3,5,0,1,50,50,50,50);