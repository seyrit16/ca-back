# chess-adventures-backend

## 1. Цель проекта

`chess-adventures-backend` — серверная часть пошаговой PvP-игры «Chess Adventures» на Go.
Сервер закрывает полный цикл партии: регистрация и JWT-аутентификация, подбор соперника,
создание партии со случайно генерируемой гексагональной картой, драфт пяти персонажей
(«snake draft»), realtime-игровая сессия по WebSocket (перемещение юнитов, рубка деревьев,
вытягивание карт, вскрытие секретов, предметы и хранилища) и SSE-уведомление о найденном
сопернике. Активная партия считается в памяти процесса (комната на игру), PostgreSQL хранит
пользователей, сессии, партии, клетки, юнитов, предметы и шаблоны.

## 2. Стек и запуск

| Что | Значение |
|---|---|
| Язык / модуль | Go 1.26, модуль `ca` |
| HTTP | Gin (`gin.Default()`), слушает `cfg.Host` = `localhost:8080` |
| WebSocket | `gorilla/websocket` |
| БД | PostgreSQL через `pgx/v5` + `pgxpool`; строки в структуры — `scany/v2/pgxscan` |
| Аутентификация | JWT HS256 (`pkg/jwtutil`), bcrypt для паролей, refresh-сессии в БД |
| Конфиг | YAML `config/config.yml` (`config.Load`) |
| Миграции | SQL в `migration/`, golang-migrate через `script/make.ps1` |
| Прочее | Redis в конфиге предусмотрен, но не используется; тестов нет |

Запуск: `go run ./cmd/app` или `script/make.ps1 run` (есть также `build`, `dev`, `migrate-up`,
`test`, `fmt`, `vet`, `lint`).

## 3. Структура папок

```
cmd/app/                         точка входа: читает конфиг, создаёт pool, репозитории,
                                 сервисы, rest-сервер, запускает Hub и ListenAndServe
config/config.yml                host, postgres url, JWT-секреты и TTL
internal/config/                 загрузка YAML в Config
internal/model/                  доменные сущности с db-тегами: User, Game, GamePlayer,
                                 GameUnit, GameCell, GameItem, ItemTemplate,
                                 CharacterTemplate, RefreshSession
internal/model/enum/             перечисления: роли, статусы партии, биомы, контент клеток,
                                 бедствия, типы предметов и их id
internal/mapgen/                 генератор гексагональной карты (radius/spawnRadius,
                                 сид-рандом, биомы, деревья, секреты, монстры, ToHTML)
internal/repository/postgres/    интерфейсы репозиториев + pgx-реализации и SQL
internal/service/                бизнес-логика: auth, user, game, game_session, game_cell,
                                 game_player, game_unit (драфт), game_item, matchmaking
internal/transport/rest/         Gin-роутер (api.go), хендлеры auth/user/matchmaking,
                                 middleware (CORS, JWT)
internal/transport/dto/          request/response-структуры и функции ToXDTO
internal/transport/websocket/    общий Hub (клиенты по комнатам), Client, конверт сообщений
internal/transport/websocket/game/   игровая комната: handler, hub, room (действия),
                                     message (payload-структуры), примеры JSON-снимков
internal/transport/websocket/draft/  драфт: handler, hub, таймеры хода и автопик
internal/room/                   in-memory состояние партии: State, DeckState, SecretReward
internal/web/auth.html           тестовая страница логина/регистрации
migration/                       up/down SQL-миграции (games, players, users, sessions,
                                 units, cells, item_templates, game_items, character_templates)
pkg/jwtutil/                     генерация/парсинг access и refresh токенов
script/make.ps1                  build/run/test/migrate/lint/clean
.idea/, app.exe                  артефакты IDE и сборки — не часть кода
```

## 4. Ключевые модули

| Модуль | Ответственность |
|---|---|
| `cmd/app/main.go` | Ручная сборка зависимостей (DI без контейнера): pool → репозитории → сервисы → `rest.New` |
| `transport/rest/api.go` | Все маршруты, создание `DraftHub` и `GameHub`, запуск `gameHub.Run` |
| `transport/rest/middleware.go` | Логирование, CORS (`http://localhost:5173`), JWT-проверка (Bearer или `?token=`), `RequireRole` |
| `transport/rest/auth.go` | login / refresh / logout, работа с refresh-cookie (сейчас закомментирована) |
| `transport/rest/matchmaking.go` | join / leave / status (SSE) |
| `service/auth.go` | Проверка пароля, выпуск пары токенов, ротация refresh, хэш SHA-256 в БД |
| `service/matchmaking.go` | Один слот ожидания в памяти, создание партии при паре, каналы подписки |
| `service/game.go` | `CreateGameWithPlayer`: запись в БД + асинхронная генерация карты |
| `service/game_session.go` | Сборка полного снимка партии из БД (game + cells + players + units + шаблоны) |
| `service/game_unit.go` | Драфт в памяти: порядок 0-1-1-0-0-1-1-0-0-1, пики, автопик, вставка юнитов из `character_templates` |
| `service/game_cell.go` | Преобразование результата mapgen в строки `game_cells` |
| `room/state.go` | Живое состояние: `model.Game`, шаблоны предметов, фазы хода, колода, награды секретов, генерация id предметов |
| `room/deck.go` | Колода карт (4 масти + скорость + 2 джокера), выдача карт, счётчики вытягиваний |
| `room/secret.go` | Взвешенный пул наград за секрет и кэш по `requestId` |
| `websocket/game/*` | Комната на игру: загрузка состояния, welcome, снимок, обработчики действий, рассылка патчей |
| `websocket/draft/*` | Драфт-сокет: welcome, `draft_state`, таймер 30 с, автопик, `draft_full`, обратный отсчёт |
| `dto/*` | Контракт транспорта: snake_case JSON для REST, мапперы в DTO, отдельные payload-структуры WS |

## 5. Как данные движутся через систему

### 5.1 Общая схема

```
HTTP/WS клиент
   │
   ▼
Gin router ── middleware (CORS, лог, JWT) ── handler
   │                                            │
   │                                            ▼
   │                                         service  (правила, координация)
   │                                            │
   │                                            ▼
   │                                    repository/postgres (pgx, SQL)
   │                                            │
   │                                            ▼
   │                                        PostgreSQL
   │
   └── WebSocket ── hub (clients[roomId]) ── room (in-memory State)
                                                │
                                                └── broadcast JSON-патчей всем в комнате
```

Ключевая особенность: **во время партии источник истины — память**. Комната один раз
загружает партию из БД (`GetGameSession`) и дальше меняет только `room.State`.

### 5.2 Аутентификация

1. `POST /api/auth/register` → `UserService.Register` → bcrypt-хэш → `INSERT users`.
2. `POST /api/auth/login` → `AuthService.Login` → `SELECT users` → bcrypt-сравнение →
   выпуск access (15 мин) и refresh (168 ч) JWT → SHA-256 refresh-токена пишется в
   `refresh_sessions` (одна сессия на пользователя, upsert по `user_id`).
3. Клиент шлёт `Authorization: Bearer <access>`; middleware парсит токен и кладёт в контекст
   `user_id`, `username`, `role`.
4. `POST /api/auth/refresh` → проверка подписи refresh → поиск хэша в БД → проверка срока →
   удаление старой сессии → выпуск новой пары (ротация).
5. `POST /api/auth/logout` → удаление сессии по refresh-cookie (cookie сейчас не выставляется,
   см. «Известные особенности»).

### 5.3 Матчмейкинг и создание партии

1. Игрок A: `POST /api/matchmaking/join` → слот `waiting` занят, ответ `{"status":"waiting"}`.
2. Игрок A открывает `GET /api/matchmaking/status` (SSE) и ждёт событие.
3. Игрок B: `join` → сервис видит ожидающего, создаёт партию
   (`games` + две строки `game_players`: `player_0`/`player_1`, seat 0/1), отправляет
   `game_id` в канал A, возвращает `{"status":"found","game_id":...}`.
4. `GameService.CreateGameWithPlayer` в горутине генерирует карту (`mapgen`, radius 20,
   spawnRadius 3) и одним батч-INSERT пишет все клетки в `game_cells`.

### 5.4 Драфт персонажей

1. Оба игрока подключаются к `WS /ws/draft/{game_id}`; каждому приходит `welcome` с
   `player_id`, затем `JoinDraft`.
2. Когда в драфте оба игрока, порядок ходов случайно перемешивается, включается
   `draft_state` и 30-секундный таймер; по таймауту — автопик шаблона `0`.
3. Клиент шлёт `{"type":"pick","payload":{"template_id":N}}`; сервис проверяет очередь и
   лимит 5 пиков, вставляет юнита в `game_units` из `character_templates` со стартовой
   позицией (своя пятёрка координат для `player_0` и `player_1`).
4. После 10 пиков — `draft_full`, 10-секундный обратный отсчёт и `draft_complete`.
5. Состояние драфта живёт в `GameUnitService.drafts` (map в памяти) и теряется при рестарте.

### 5.5 Игровая сессия

1. `WS /ws/game/{game_id}` → JWT → `GameWsHandler` → `GameHub`.
2. Hub держит `rooms map[string]*Room`. Первый клиент создаёт комнату: `GetGameSession`
   (game + board + players + units) + `GetItemTemplates` → `room.NewState` → `room.Run`.
3. `AddPlayer` находит игрока по `user_id`, запоминает `player_id`, отправляет `welcome`.
   Когда в комнате двое (или игрок переподключился), всем уходит `state_snapshot` с `GameDTO`.
4. Дальше клиент шлёт «плоские» кадры (`type` + поля действия), комната меняет `State` и
   рассылает всем `state_patch` (плоский JSON с `type`/`entity` и полями payload).
5. При отключении клиент удаляется из комнаты; пустая комната удаляется из hub.

### 5.6 Колода, секреты, предметы

- `DeckState`: 67 карт (4 масти × 13 значений, 13 карт скорости, 2 джокера). Вытягивание
  уменьшает счётчик `DrawCounts[playerId][unitId]`; при исчерпании колода перемешивается заново.
- `SecretReward`: взвешенный пул наград; `take_chest` кладёт выбранную награду в кэш по
  `requestId`, `take_chest_resource` превращает её в предмет.
- `State.SetupItemByTemplate` создаёт `GameItem` с локальным id и кладёт в `Storage` игрока
  или `Pocket` юнита в зависимости от `default_location` шаблона. В БД предметы не пишутся
  (TODO), хотя репозиторий и сервис для этого уже есть.

### 5.7 Генерация карты

`mapgen.MapGenerator` строит гекс-поле (radius 20 → 1261 клетка), расставляет воду и горы
кластерами, деревья, секреты и монстров, исключая красные зоны у точек спавна. Результат —
`Tiles`, `Coords`, `Items`; `GameCellService` конвертирует их в `model.GameCell` и сохраняет.
Тот же генератор используется в отладочном `GET /test/map/generate` (отдаёт HTML).

## 6. Важные соглашения

**Слои и зависимости.** `transport → service → repository → БД`; `model` и `enum` не зависят
ни от чего, кроме стандартной библиотеки и uuid. Репозитории объявлены интерфейсами в пакете
`postgres`, реализации — неэкспортируемые структуры, конструкторы `NewXRepository` возвращают
интерфейс (у `gameItemRepository`/`itemTemplateRepository` — конкретный указатель).

**Именование.** Конструкторы — `NewX`; сервисы — `XService`, репозитории — `XRepository`,
хендлеры — `XHandler`; DTO-мапперы — `ToXDTO` / `XToDTO`. В Go-полях используется `Id`
(`UserId`, `GameId`, `PlayerId`), не `ID`. Пакет `ca` — короткий путь импорта (`ca/internal/...`).

**Модели и БД.** Таблицы и колонки — `snake_case`, теги `db:"..."`; составные первичные ключи
(`game_id, player_id`, `game_id, player_id, unit_id`, `game_id, q, r`); UUID — `gen_random_uuid()`.
Чтение — `pgx.RowToStructByNameLax` или `pgxscan`; ошибки репозиториев оборачиваются префиксом
`[ИмяФункции]`.

**Транспорт.** REST-DTO — `snake_case` (`access_token`, `rating_delta`, `player_id`).

**Идентификаторы.** `player_0` / `player_1` — id игрока внутри партии (последний символ — seat);
`unit_0_0` … `unit_1_4` — `<seat игрока>_<номер юнита>`; ключ клетки на доске — строка `"q,r"`;
в JWT и БД — UUID.

**Ошибки.** REST: `{"error": "..."}`; коды 400/401/403/500.
WS: `{"type":"error","entity":"","payload":{"message":"..."}}`. Логи — `log`/`slog`,
фатальные ошибки старта — `log.Fatal`.

**Состояние.** Всё игровое состояние — в памяти (`room.State`, `GameUnitService.drafts`,
`MatchmakingService.waiting`); горутины комнаты и hub синхронизируются мьютексом/каналами.
`room` — отдельный пакет, чтобы не тянуть транспорт в игровую логику.

**Конфиг.** Секреты и TTL — только в `config/config.yml`; TTL парсятся YAML в `time.Duration`
(`15m`, `168h`).

