# API — chess-adventures-backend

Документация описывает: HTTP-маршруты Gin,
SSE-эндпоинт матчмейкинга и два WebSocket-канала (драфт и игровая сессия).
Базовый адрес по умолчанию — `http://localhost:8080` (порт задаётся `host` в `config/config.yml`).

## Содержание

- [Общие сведения](#общие-сведения)
- [Сводная таблица маршрутов](#сводная-таблица-маршрутов)
- [REST: служебные маршруты](#rest-служебные-маршруты)
- [REST: аутентификация](#rest-аутентификация)
- [REST: матчмейкинг](#rest-матчмейкинг)
- [WebSocket: игровая сессия](#websocket-игровая-сессия)
- [WebSocket: драфт](#websocket-драфт)
- [Приложение A. Перечисления](#приложение-a-перечисления)
- [Приложение B. Шаблоны](#приложение-b-шаблоны)
- [Приложение C. Тексты ошибок](#приложение-c-тексты-ошибок)

## Общие сведения

| Параметр | Значение |
|---|---|
| Протокол | HTTP/1.1 + WebSocket (gorilla/websocket) |
| Формат | JSON (`application/json; charset=utf-8`), кроме HTML-страниц и SSE |
| Авторизация | `Authorization: Bearer <access_token>`; для WS и SSE допустим `?token=<access_token>` |
| JWT | HS256, claims: `UserId`, `Username`, `Role`, `exp`, `iat`; access — 15 мин, refresh — 168 ч |
| CORS | `Access-Control-Allow-Origin: http://localhost:5173`, credentials, методы `POST, GET, OPTIONS, PUT, DELETE`, заголовки `Content-Type, Authorization` |
| Preflight | любой `OPTIONS` → `204 No Content` |

Каждый запрос логируется в stdout: метод, путь и все заголовки.

### Форматы ошибок

- Большинство REST-хендлеров: `{"error": "<текст>"}`.
- `POST /api/auth/register` при ошибке биндинга/БД отдаёт только статус (тело пустое,
  ошибка пишется через `AbortWithError`).
- `GET /test/game/{game_id}` при ошибке сервиса отдаёт JSON-строку: `"<текст ошибки>"`.
- WebSocket: `{"type":"error","entity":"","payload":{"message":"<текст>"}}`.

### Идентификаторы

- `player_id` — `player_0` или `player_1` (последний символ — seat).
- `unit_id` — `unit_<seat игрока>_<номер юнита>`, например `unit_0_3`.
- Ключ клетки доски — строка `"q,r"`, например `"-3,5"`.
- `game_id`, `user_id` — UUID.

## Сводная таблица маршрутов

| Метод | Путь | Авторизация | Назначение |
|---|---|---|---|
| GET | `/test/ping` | нет | Проверка живости |
| GET | `/test/map/generate` | нет | HTML-визуализация сгенерированной карты |
| GET | `/test/game/{game_id}` | нет | Отладочный дамп партии из БД |
| GET | `/auth` | нет | HTML-страница логина/регистрации |
| POST | `/api/auth/register` | нет | Регистрация |
| POST | `/api/auth/login` | нет | Вход, выдача токенов |
| POST | `/api/auth/refresh` | нет (нужен refresh-токен в теле) | Ротация токенов |
| POST | `/api/auth/logout` | нет | Выход (по refresh-cookie) |
| POST | `/api/matchmaking/join` | JWT | Встать в очередь / создать пару |
| DELETE | `/api/matchmaking/leave` | JWT | Покинуть очередь |
| GET | `/api/matchmaking/status` | JWT | SSE-ожидание соперника |
| GET | `/ws/game/{game_id}` | JWT | WebSocket игровой сессии |
| GET | `/ws/draft/{game_id}` | JWT | WebSocket драфта персонажей |

## REST: служебные маршруты

### GET /test/ping

Проверка доступности сервера.

- Параметров нет.
- Ответ `200 OK`: JSON-строка.

```bash
curl http://localhost:8080/test/ping
```

```json
"pong"
```

### GET /test/map/generate

Создаёт `mapgen.NewMapGenerator(20, 3)`, генерирует карту и возвращает её HTML-представление
(`generator.ToHTML()`). Каждый вызов генерирует новую карту.

- Параметров нет.
- Ответ `200 OK`, `Content-Type: text/html; charset=utf-8` — HTML-документ с картой.

```bash
curl http://localhost:8080/test/map/generate -o map.html
```

### GET /test/game/{game_id}

Отладочный дамп партии: `GameSessionService.GetGameSession` загружает игру, клетки, игроков
и юнитов. Возвращается `model.Game` **без json-тегов**, поэтому ключи — имена Go-полей
(`Id`, `OwnerId`, `Board`, …), а не snake_case. Вложенные `Storage`/`Pocket` не загружаются.

Параметры пути:

| Имя | Тип | Описание |
|---|---|---|
| `game_id` | UUID (string) | Идентификатор партии |

Ответы:

| Код | Тело |
|---|---|
| 200 | `model.Game` в JSON |
| 400 | `{"error":"invalid game_id"}` — невалидный UUID |
| 400 | JSON-строка с текстом ошибки загрузки, например `"[LoadGame]\n..."` |

```bash
curl http://localhost:8080/test/game/66381913-c151-4cf0-89db-05aef255b24e
```

```json
{
  "Id": "66381913-c151-4cf0-89db-05aef255b24e",
  "OwnerId": null,
  "Status": 0,
  "WinnerId": "00000000-0000-0000-0000-000000000000",
  "RatingDelta": 0,
  "BoardRadius": 0,
  "Seeds": "",
  "TurnCount": 0,
  "CurrentTurnPlayerId": "",
  "NumberCharacters": 5,
  "NumberCubes": 5,
  "CurrentDisaster": 0,
  "DisasterStartedAt": 0,
  "DisasterZone": 0,
  "CreatedAt": "2026-10-08T13:23:06.498501+03:00",
  "UpdatedAt": "2026-10-08T13:23:06.498501+03:00",
  "FinishedAt": null,
  "Board": {
    "0,0": {
      "GameId": "66381913-c151-4cf0-89db-05aef255b24e",
      "Q": 0,
      "R": 0,
      "GenerationGroup": null,
      "IndexInGroup": null,
      "Biome": 0,
      "Content": 0,
      "RadiusContent": null
    }
  },
  "Players": {
    "player_0": {
      "GameId": "66381913-c151-4cf0-89db-05aef255b24e",
      "PlayerId": "player_0",
      "UserId": "e89b8faa-4d02-44e3-a95a-5a2162923035",
      "Seat": 0,
      "IsWinner": false,
      "Units": {},
      "Storage": null
    }
  }
}
```

### GET /auth

Отдаёт `internal/web/auth.html` — тестовую страницу с формами входа и регистрации.
Страница обращается к `http://localhost:8080` и после логина перебрасывает токены на
`http://localhost:7777/callback`.

```bash
curl http://localhost:8080/auth
```

## REST: аутентификация

### POST /api/auth/register

Регистрация пользователя. Пароль хэшируется bcrypt, `role` берётся из DEFAULT БД (`user`).

Тело запроса (`application/json`):

| Поле | Тип | Обязательное | Описание |
|---|---|---|---|
| `nickname` | string | нет | Отображаемое имя |
| `username` | string | нет | Логин, UNIQUE |
| `password` | string | нет | Пароль (хэшируется) |
| `email` | string | нет | Email, UNIQUE |
| `color` | string | нет | Принимается DTO, но **не сохраняется** |

Валидации полей нет; пустые значения уйдут в БД и упрутся в ограничения.

Ответы:

| Код | Тело | Причина |
|---|---|---|
| 200 | `{}` | Пользователь создан |
| 400 | пусто | Ошибка биндинга JSON |
| 500 | пусто | Ошибка БД (например, дубликат username/email) |

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"nickname":"Игрок","username":"player1","password":"secret","email":"p1@example.com"}'
```

```json
{}
```

### POST /api/auth/login

Проверяет пару логин/пароль и выдаёт access + refresh токены. Хэш refresh-токена
(SHA-256) сохраняется в `refresh_sessions`; на пользователя — одна сессия (upsert).

Тело запроса:

| Поле | Тип | Описание |
|---|---|---|
| `username` | string | Логин |
| `password` | string | Пароль |

Успешный ответ `200 OK`:

| Поле | Тип | Описание |
|---|---|---|
| `access_token` | string | JWT для `Authorization: Bearer` |
| `refresh_token` | string | JWT для `/api/auth/refresh` |
| `expires_in` | int64 | Всегда `900` секунд (15 мин, захардкожено) |

Ответы:

| Код | Тело |
|---|---|
| 200 | `{"access_token":"…","refresh_token":"…","expires_in":900}` |
| 400 | `{"error":"<ошибка биндинга>"}` |
| 401 | `{"error":"invalid credentials"}` |

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"player1","password":"secret"}'
```

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9…",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9…",
  "expires_in": 900
}
```

### POST /api/auth/refresh

Ротация токенов: проверяется подпись refresh-токена, наличие его хэша в БД и срок действия,
старая сессия удаляется, выпускается новая пара.

Тело запроса:

| Поле | Тип | Описание |
|---|---|---|
| `refresh_token` | string | Действующий refresh-токен |

Ответы:

| Код | Тело |
|---|---|
| 200 | `{"access_token":"…","refresh_token":"…","expires_in":900}` |
| 400 | `{"error":"<ошибка биндинга>"}` |
| 401 | `{"error":"invalid refresh token"}` / `{"error":"session not found"}` / `{"error":"refresh token expired"}` |

При ошибке дополнительно очищается cookie `refresh_token` (path `/api/auth`).

```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9…"}'
```

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9…",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9…",
  "expires_in": 900
}
```

### POST /api/auth/logout

Удаляет refresh-сессию, если в запросе есть cookie `refresh_token`, и очищает эту cookie.

- Тело не читается; JWT не требуется.
- Cookie: HttpOnly, Secure, path `/api/auth`.
- Важно: `login` и `refresh` cookie **не выставляют** (код установки закомментирован),
  поэтому на практике logout удаляет сессию только если cookie появилась иным способом.

Ответ `200 OK`:

```json
{"message":"logged out"}
```

```bash
curl -X POST http://localhost:8080/api/auth/logout \
  -b "refresh_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9…"
```

## REST: матчмейкинг

Все три маршрута требуют JWT. Подпись access-токена проверяется middleware; при ошибке —
`401 {"error":"unauthorized"}` (токена нет) или `401 {"error":"invalid or expired token"}`.
Токен берётся из `Authorization: Bearer …`, а если заголовка нет — из query-параметра `token`.

### POST /api/matchmaking/join

Ставит игрока в очередь. В сервисе один слот ожидания: если там уже другой игрок,
создаётся партия, и оба получают её id; иначе текущий игрок становится ожидающим.

- Параметров и тела нет.
- Ответ `200 OK`:
  - `{"status":"waiting"}` — игрок поставлен в очередь;
  - `{"status":"found","game_id":"<uuid>"}` — пара найдена, партия создана.
- `500 {"error":"<текст>"}` — ошибка создания партии в БД.

```bash
curl -X POST http://localhost:8080/api/matchmaking/join \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

```json
{"status":"found","game_id":"66381913-c151-4cf0-89db-05aef255b24e"}
```

### DELETE /api/matchmaking/leave

Убирает игрока из очереди и удаляет его канал подписки.

- Параметров нет.
- Ответ `200 OK`: `{"status":"left"}`.

```bash
curl -X DELETE http://localhost:8080/api/matchmaking/leave \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

### GET /api/matchmaking/status (SSE)

Ожидание соперника через Server-Sent Events. Хендлер подписывается на канал игрока,
отдаёт одно событие и завершает ответ.

Заголовки ответа: `Content-Type: text/event-stream`, `Cache-Control: no-cache`,
`Connection: keep-alive`.

Событие:

```
event:match_found
data:{"game_id":"66381913-c151-4cf0-89db-05aef255b24e"}
```

- Если соединение оборвано клиентом до матча, игрок удаляется из очереди (`Leave`).
- Если `join` не вызывался, канал не сохраняется в сервисе, и событие не придёт.
- Для EventSource токен передаётся в query: `?token=<access_token>`.

```bash
curl -N "http://localhost:8080/api/matchmaking/status?token=$ACCESS_TOKEN"
```

## WebSocket: игровая сессия

### Подключение

```
ws://localhost:8080/ws/game/{game_id}?token=<access_token>
```

- `game_id` — UUID существующей партии.
- Авторизация — как у REST: `Authorization` не обязателен, для браузера используйте `?token=`.
  Без валидного токена — HTTP `401` до апгрейда.
- `CheckOrigin` разрешает любой Origin.
- В комнате максимум 2 клиента; третий получает ошибку `already full`.

### Формат кадров

Входящие кадры — объект `IncomingMessage`:

| Поле | Тип | Описание |
|---|---|---|
| `type` | string | Тип действия (или `"action"`) |
| `actionType` | string | Альтернатива `type`, если `type == "action"` |
| `payload` | object | Используется только действием `card_to_unit` |

Для всех действий, кроме `card_to_unit`, сервер парсит **весь кадр** как структуру действия,
поэтому поля действия лежат на верхнем уровне, рядом с `type`, а не внутри `payload`.

Исходящие кадры бывают двух видов:

1. Конверт: `{"type":"…","entity":"…","payload":{…}}` — для `welcome`, `state_snapshot`, `error`.
2. «Плоский» патч: `{"type":"state_patch","entity":"<сущность>","<поля payload>":…}` — для
   игровых изменений (`sendAllFlat` поднимает поля на верхний уровень).

### Сценарий подключения

1. Клиент открывает сокет.
2. Сервер регистрирует клиента в комнате; при первом клиенте комната загружает партию из БД.
3. Клиенту приходит `welcome` (если его `user_id` найден среди игроков партии).
4. Когда в комнате двое (или игрок переподключился), всем уходит `state_snapshot`.
5. Далее — действия клиента и патчи сервера.

### Исходящее: welcome

```json
{
  "type": "welcome",
  "entity": "",
  "payload": {"type": "", "player_id": "player_0", "seat": 0}
}
```

`payload.type` в игровом канале присутствует, но всегда пустой. Если пользователь не
является игроком этой партии, `welcome` не отправляется, но соединение остаётся.

### Исходящее: state_snapshot

Полный снимок партии (`dto.GameDTO`) в конверте. Отправляется, когда комната заполнена,
и при повторном подключении уже присутствующего игрока.

```json
{
  "type": "state_snapshot",
  "entity": "",
  "payload": {
    "status": 0,
    "rating_delta": 0,
    "board_radius": 0,
    "seeds": "",
    "turn_count": 0,
    "current_turn_player_id": "",
    "number_characters": 5,
    "number_cubes": 5,
    "current_disaster": 0,
    "disaster_started_at": 0,
    "disaster_zone": 0,
    "created_at": "2026-10-08T13:23:06.498501+03:00",
    "updated_at": "2026-10-08T13:23:06.498501+03:00",
    "finished_at": null,
    "item_templates": [
      {"id": 3, "itemName": "Heal", "description": "Вылечите одного своего персонажа", "defaultLocation": 0}
    ],
    "board": {
      "-1,-1": {
        "q": -1, "r": -1,
        "generation_group": 12, "index_in_group": 77,
        "biome": 0, "content": 0, "radius_content": 0
      }
    },
    "players": {
      "player_0": {
        "player_id": "player_0",
        "seat": 0,
        "units": {
          "unit_0_0": {
            "player_id": "player_0", "unit_id": "unit_0_0", "seat": 0,
            "position_q": -18, "position_r": 18,
            "character_template_id": 0, "name": "Аргул",
            "hp": 1, "max_hp": 1, "attack": 5, "defense": 0,
            "luck": 0, "speed": 1, "critical_chance": 50, "critical_damage": 50,
            "accuracy": 50, "dodge": 50,
            "equip_head": null, "equip_body": null, "equip_boots": null, "equip_weapon": null,
            "is_alive": true
          }
        }
      }
    }
  }
}
```

Поля `payload`:

| Поле | Тип | Описание |
|---|---|---|
| `status` | int | Статус партии (см. приложение A) |
| `rating_delta` | int16 | Изменение рейтинга |
| `board_radius` | int16 | Радиус доски |
| `seeds` | string | Сид генерации |
| `turn_count` | int | Номер хода |
| `current_turn_player_id` | string | Чей ход (`player_0`/`player_1`) |
| `number_characters` | int | Персонажей на игрока |
| `number_cubes` | int | Кубиков |
| `current_disaster` | int | Текущее бедствие |
| `disaster_started_at` | int | Время начала бедствия |
| `disaster_zone` | int16 | Зона бедствия |
| `created_at`, `updated_at` | string (RFC 3339) | Метки времени |
| `finished_at` | string / null | Время завершения |
| `item_templates` | array | Шаблоны предметов: `id`, `itemName`, `description`, `defaultLocation` |
| `board` | object | Ключ `"q,r"` → клетка |
| `players` | object | Ключ `player_id` → игрок с `units` |

Клетка: `q`, `r` (int16), `generation_group`, `index_in_group`, `radius_content`
(int16 или null), `biome`, `content` (int).

### Исходящее: state_patch

Патч состояния конкретной сущности. Поля действия лежат на верхнем уровне.

| `entity` | Дополнительные поля | Когда приходит |
|---|---|---|
| `unit_move` | `request_id`, `unitId`, `path`, `diceIndex` | После действия `move` |
| `dice_roll` | `requestId`, `id` | После `dice_roll` |
| `dice_thrown_flag` | `id`, `diceValues` | После `dice_thrown_flag` |
| `tree_chop` | `requestId`, `unitId`, `path`, `diceIndex`, `id` | После `tree_chop` |
| `turn_phase` | `id`, `phase`, `draws` | Смена фазы на `draw_cards` |
| `draw_card` | `id`, `unitId`, `card`, `draws` | После `draw_card` |
| `take_chest` | `id`, `unitId`, `requestId`, `reward` | После `take_chest` |
| `player_storage` | `id`, `items` | Изменение склада игрока |
| `unit_item` | `id`, `unitId`, `items` | Изменение кармана юнита |
| `unit_update` | вложенный `payload`: `player_id`, `unit_id`, `hp`, `max_hp`, `attack`, `defense`, `luck`, `speed` | После `card_to_unit` (единственный вложенный патч) |

`card` — `{suit, baseRank, bonusValue, unitId}`.
`draws` — `map[string]int8`: `unit_id → сколько карт осталось вытянуть` (ключ `""` — карты игрока).
`items` — массив `{id, templateId, location, unitId}`.

### Входящие действия

Общий вид: `{"type":"<action>", ...поля действия}`.
Поддерживается и вариант `{"type":"action","actionType":"<action>", ...поля}`.

| `type` | Назначение |
|---|---|
| `move` | Перемещение юнита по пути |
| `dice_roll` | Запрос броска кубиков |
| `dice_thrown_flag` | Фиксация броска, начисление карт за дубли |
| `tree_chop` | Рубка дерева |
| `end_turn_phase` | Переход к следующей фазе хода |
| `draw_card` | Вытянуть карту |
| `card_to_unit` | Применить карту к юниту |
| `take_chest` | Вскрыть секрет |
| `take_chest_resource` | Забрать награду секрета |

Неизвестный `type` → `error` с текстом `unknown action`.

#### move

Перемещает юнита в последнюю точку пути. Промежуточные клетки не проверяются (TODO),
очередь хода не проверяется.

| Поле | Тип | Описание |
|---|---|---|
| `request_id` | string | Клиентский id запроса |
| `unitId` | string | Id юнита |
| `path` | array `{q,r}` | Путь; берётся последняя точка |
| `diceIndex` | int | Индекс кубика |

```json
{"type":"move","request_id":"req-1","unitId":"unit_0_1","path":[{"q":-18,"r":18},{"q":-17,"r":17}],"diceIndex":8}
```

Ответ всем:

```json
{"type":"state_patch","entity":"unit_move","request_id":"req-1","unitId":"unit_0_1","path":[{"q":-18,"r":18},{"q":-17,"r":17}],"diceIndex":8}
```

#### dice_roll

| Поле | Тип | Описание |
|---|---|---|
| `requestId` | string | Клиентский id запроса |
| `id` | string | Игнорируется, сервер подставляет `player_id` |

```json
{"type":"dice_roll","requestId":"req-2"}
```

```json
{"type":"state_patch","entity":"dice_roll","requestId":"req-2","id":"player_0"}
```

#### dice_thrown_flag

Сервер сам бросает 5 кубиков d6, считает дубли и начисляет столько карт игроку
(ключ `""` в счётчиках).

| Поле | Тип | Описание |
|---|---|---|
| `id` | string | Игнорируется |
| `diceValues` | int[] | Игнорируется, значения генерирует сервер |

```json
{"type":"dice_thrown_flag"}
```

```json
{"type":"state_patch","entity":"dice_thrown_flag","id":"player_0","diceValues":[3,3,3,1,6]}
```

#### tree_chop

Очищает клетку от дерева, начисляет 1 карту юниту и создаёт предмет «дерево»
(шаблон `0`) на складе игрока.

| Поле | Тип | Описание |
|---|---|---|
| `requestId` | string | Клиентский id |
| `unitId` | string | Юнит, который рубит |
| `path` | array `{q,r}` | Путь; последняя точка — клетка с деревом |
| `diceIndex` | int | Индекс кубика |
| `id` | string | Игнорируется, в ответе заменяется на `unitId` |

```json
{"type":"tree_chop","requestId":"req-3","unitId":"unit_0_2","path":[{"q":-15,"r":15}],"diceIndex":2}
```

Патчи: `tree_chop`, затем `player_storage`:

```json
{"type":"state_patch","entity":"tree_chop","requestId":"req-3","unitId":"unit_0_2","path":[{"q":-15,"r":15}],"diceIndex":2,"id":"unit_0_2"}
{"type":"state_patch","entity":"player_storage","id":"player_0","items":[{"id":0,"templateId":0,"location":0}]}
```

#### end_turn_phase

Переключает фазу на следующую по порядку `character_moves → draw_cards → reveal_secrets`.
Для `draw_cards` рассылается `turn_phase` со счётчиками вытягиваний; для `reveal_secrets`
ничего не рассылается.

| Поле | Тип | Описание |
|---|---|---|
| `phase` | string | Текущая фаза: `character_moves`, `draw_cards`, `reveal_secrets` |

```json
{"type":"end_turn_phase","phase":"character_moves"}
```

```json
{"type":"state_patch","entity":"turn_phase","id":"player_0","phase":"draw_cards","draws":{"unit_0_1":1,"":2}}
```

#### draw_card

Вытягивает карту из колоды, если у игрока/юнита есть неиспользованные вытягивания.
Джокеры (suit 5 и 6) дополнительно создают предмет `10`/`11` на складе игрока.

| Поле | Тип | Описание |
|---|---|---|
| `unitId` | string | Юнит, которому предназначена карта; `""` — карта игрока |

```json
{"type":"draw_card","unitId":""}
```

```json
{
  "type":"state_patch","entity":"draw_card","id":"player_0","unitId":"",
  "card":{"suit":0,"baseRank":7,"bonusValue":0,"unitId":""},
  "draws":{"unit_0_1":1}
}
```

#### card_to_unit

Единственное действие, читающее вложенный `payload`. Меняет характеристику юнита:
`change = bonusValue + baseRank` при `isPositive: true`, иначе `bonusValue - baseRank`.

| Поле | Тип | Описание |
|---|---|---|
| `suit` | int | 0 — HP, 1 — defense, 2 — attack, 3 — luck, 4 — speed |
| `baseRank` | int | Номинал карты |
| `bonusValue` | int | Бонус |
| `unitId` | string | Целевой юнит |
| `isPositive` | bool | Знак изменения |

```json
{"type":"card_to_unit","payload":{"suit":2,"baseRank":7,"bonusValue":0,"unitId":"unit_0_1","isPositive":true}}
```

```json
{"type":"state_patch","entity":"unit_update","payload":{"player_id":"player_0","unit_id":"unit_0_1","attack":12}}
```

#### take_chest

Взвешенно выбирает награду, кладёт её в кэш по `requestId` и сообщает id шаблона награды
в поле `reward`.

| Поле | Тип | Описание |
|---|---|---|
| `requestId` | string | Обязателен; ключ кэша награды |
| `unitId` | string | Юнит, который вскрывает секрет |

```json
{"type":"take_chest","requestId":"req-4","unitId":"unit_0_1"}
```

```json
{"type":"state_patch","entity":"take_chest","id":"player_0","unitId":"unit_0_1","requestId":"req-4","reward":8}
```

#### take_chest_resource

Превращает закэшированную награду в предмет. Для шаблонов `1` (Buff) и `2` (Debuff)
обработчик завершается без ответа (TODO: показ колоды).

| Поле | Тип | Описание |
|---|---|---|
| `requestId` | string | Ключ, использованный в `take_chest` |
| `unitId` | string | Юнит-получатель (для кармана) |

```json
{"type":"take_chest_resource","requestId":"req-4","unitId":"unit_0_1"}
```

Ответ — `player_storage` (награда на склад) или `unit_item` (награда в карман):

```json
{"type":"state_patch","entity":"unit_item","id":"player_0","unitId":"unit_0_1","items":[{"id":1,"templateId":8,"location":1,"unitId":"unit_0_1"}]}
```

### Исходящее: error

```json
{"type":"error","entity":"","payload":{"message":"unknown action"}}
```

Тексты ошибок — в приложении C.

## WebSocket: драфт

### Подключение

```
ws://localhost:8080/ws/draft/{game_id}?token=<access_token>
```

- Требуется JWT (как в игровом канале).
- Пользователь должен быть участником партии (`game_players`); иначе соединение закрывается
  без сообщения об ошибке.
- При подключении сервер:
  1. регистрирует клиента в `DraftHub`;
  2. отправляет `welcome` с `seat: -1`;
  3. вызывает `JoinDraft` (первый пик создаёт состояние драфта);
  4. если в драфте уже двое и очередь определена — рассылает `draft_state` и запускает
     таймер хода 30 секунд.

### Исходящее: welcome

В драфте `welcome` отправляется «плоским» кадром, без конверта `entity`/`payload`:

```json
{"type":"welcome","player_id":"player_0","seat":-1}
```

### Исходящее: draft_state

```json
{
  "type": "draft_state",
  "entity": "",
  "payload": {
    "your_id": "",
    "current_turn": "player_0",
    "turn_index": 0,
    "picks": {"player_0": [], "player_1": []}
  }
}
```

| Поле | Тип | Описание |
|---|---|---|
| `your_id` | string | Всегда пустая строка (не заполняется сервером) |
| `current_turn` | string | Чей пик |
| `turn_index` | int | 0…10, номер пика в общей последовательности |
| `picks` | object | `player_id → []template_id` |

Порядок пиков — snake draft: `player_0, player_1, player_1, player_0, player_0, player_1,
player_1, player_0, player_0, player_1` (по 5 юнитов каждому).

### Входящее: pick

```json
{"type":"pick","payload":{"template_id":3}}
```

| Поле | Тип | Описание |
|---|---|---|
| `template_id` | int16 | Id шаблона персонажа (0…9) |

Ошибки (в конверте `error`): `draft not found`, `waiting for opponent`, `not your turn`,
`already picked 5 units`, `Error on pick unit`, `player not found`.

После успешного пика юнит сразу вставляется в `game_units`, рассылается новый
`draft_state`, таймер перезапускается. После 10-го пика — `draft_full`, обратный отсчёт
и `draft_complete`.

### Исходящее: countdown

```json
{"type":"countdown","entity":"","payload":{"seconds":30}}
```

Приходит каждую секунду: `30…0` на ход драфта и `10…0` перед завершением.
По истечении таймера хода сервер делает автопик шаблона `0`.

### Исходящее: draft_full и draft_complete

```json
{"type":"draft_full","entity":"","payload":null}
{"type":"draft_complete","entity":"","payload":null}
```

### Стартовые позиции юнитов

При пике юнит получает позицию по номеру места (seat 0…4):

| seat | `player_0` (q, r) | `player_1` (q, r) |
|---|---|---|
| 0 | -18, 18 | 18, -18 |
| 1 | -19, 19 | 19, -19 |
| 2 | -19, 20 | 19, -20 |
| 3 | -20, 19 | 20, -19 |
| 4 | -20, 20 | 20, -20 |

`unit_id` формируется как `unit_<последний символ player_id>_<seat>`, например `unit_0_0`.

## Приложение A. Перечисления

### GameStatus (`internal/model/enum/game_status.go`)

| Значение | Имя | Смысл |
|---|---|---|
| 0 | `waiting` | Ожидание |
| 1 | `in_progress` | Идёт |
| 2 | `finished` | Завершена |
| 3 | `abandoned` | Брошена |

### BiomeType

| Значение | Имя |
|---|---|
| 0 | `none` |
| 1 | `mountain` |
| 2 | `water` |

### ContentType

| Значение | Имя |
|---|---|
| 0 | `none` |
| 1 | `tree` |
| 2 | `secret` |
| 3 | `monster` |

### DisasterType

| Значение | Имя |
|---|---|
| 0 | calm |
| 1 | crimson night |
| 2 | freeze |
| 3 | fire |
| 4 | sandstorm |

### ItemLocationType

| Значение | Имя | Куда попадает предмет |
|---|---|---|
| 0 | `warehouse` | Склад игрока |
| 1 | `pocket` | Карман юнита |
| 2 | `equip_head` | Голова |
| 3 | `equip_body` | Тело |
| 4 | `equip_boots` | Ботинки |
| 5 | `equip_weapon` | Оружие |

### CardSuit (`internal/room/deck.go`)

| Значение | Смысл |
|---|---|
| 0 | HP |
| 1 | Defense |
| 2 | Attack |
| 3 | Luck |
| 4 | Speed |
| 5 | Red Joker (даёт предмет 10) |
| 6 | Black Joker (даёт предмет 11) |

Колода: 4 масти × значения `2…13, 15` (52 карты) + 13 карт скорости + 2 джокера = 67 карт.
При исчерпании колода собирается заново.

### Role

`user`, `admin`, `superadmin`, `tester` (хранится строкой в `users.role`, по умолчанию `user`).

## Приложение B. Шаблоны

### Предметы (`item_templates`, миграция 000007)

| id | name | default_location |
|---|---|---|
| 0 | Tree | 0 |
| 1 | Buff | 0 |
| 2 | Debuff | 0 |
| 3 | Heal | 0 |
| 4 | Buff/debuff | 0 |
| 5 | Provocation | 0 |
| 6 | Ego strike | 0 |
| 7 | Teleport | 1 |
| 8 | Camp | 1 |
| 9 | Return_stone | 1 |
| 10 | Red_joker | 0 |
| 11 | Black_joker | 0 |

### Персонажи (`character_templates`, миграция 000009)

| id | Имя | hp | max_hp | attack | defense | luck | speed |
|---|---|---|---|---|---|---|---|
| 0 | Аргул | 1 | 1 | 5 | 0 | 0 | 1 |
| 1 | Асока | 3 | 3 | 3 | 1 | 0 | 2 |
| 2 | Контроллер | 5 | 5 | 3 | 5 | 0 | 1 |
| 3 | Мегатрон | 5 | 5 | 3 | 5 | 0 | 1 |
| 4 | Нуб Сайбот | 5 | 5 | 3 | 5 | 0 | 1 |
| 5 | Сакура | 5 | 5 | 3 | 5 | 0 | 1 |
| 6 | Стенд | 5 | 5 | 3 | 5 | 0 | 1 |
| 7 | Хантер | 5 | 5 | 3 | 5 | 0 | 1 |
| 8 | Цирик | 5 | 5 | 3 | 5 | 0 | 1 |
| 9 | Эцио | 5 | 5 | 3 | 5 | 0 | 1 |

У всех: `critical_chance = critical_damage = accuracy = dodge = 50`.

## Приложение C. Тексты ошибок

### REST

| Где | Код | Текст |
|---|---|---|
| JWT middleware | 401 | `unauthorized` |
| JWT middleware | 401 | `invalid or expired token` |
| `RequireRole` (не использован в маршрутах) | 403 | `forbidden` |
| login | 401 | `invalid credentials` |
| refresh | 401 | `invalid refresh token`, `session not found`, `refresh token expired` |
| matchmaking join | 500 | текст ошибки создания партии |
| `/test/game/{id}` | 400 | `invalid game_id` |

### WebSocket (game)

| Текст | Причина |
|---|---|
| `unknown action` | Неизвестный `type`/`actionType` |
| `Invalid payload\n<err>` | Не удалось разобрать кадр действия |
| `no found playerId` | Клиенту не сопоставлен `player_id` |
| `already full` | В комнате уже 2 игрока |
| `unknown phase` | Неизвестная фаза в `end_turn_phase` |
| `no found request id` | Пустой `requestId` в `take_chest*` |
| `no found reward on cache` | Нет награды по `requestId` |
| `no found item template` | Неизвестный шаблон награды |
| `SetupItemByTemplate error` | Ошибка создания предмета |
| `DrawCard error` | Нет доступных вытягиваний |
| `Invalid suit` | Масть вне 0…4 в `card_to_unit` |

### WebSocket (draft)

| Текст | Причина |
|---|---|
| `player not found` | Пользователь не игрок партии |
| `draft not found` | Нет состояния драфта |
| `waiting for opponent` | Второй игрок ещё не подключился |
| `not your turn` | Пик не в свою очередь |
| `already picked 5 units` | Лимит пиков |
| `Error on pick unit` | Ошибка вставки юнита в БД |
