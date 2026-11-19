# scripts/make.ps1

param(
    [Parameter(Position = 0)]
    [string]$Command,
    [string]$Version
)

# ##################################################
# CONFIG
# ##################################################

$APP_NAME = 'chess-adventures'
$BUILD_DIR = './build'
$CMD_DIR = './cmd/app'
$MIGRATIONS_PATH = './migration'
$DB_URL = 'postgres://ca_user:ca_password@localhost:5432/chess-adventures-db?sslmode=disable'

# ##################################################
# FUNCTIONS
# ##################################################

function Build
{
    go build -o "$BUILD_DIR/$APP_NAME.exe" $CMD_DIR
}

function Run
{
    go run "$CMD_DIR/main.go"
}

function Dev
{
    Build
    & "$BUILD_DIR/$APP_NAME.exe"
}

function Test
{
    go test ./...
}

function TestVerbose
{
    go test -v ./...
}

function TestCover
{
    go test -cover ./...
}

function MigrateUp
{
    migrate -path $MIGRATIONS_PATH -database $DB_URL up
}

function MigrateDown
{
    migrate -path $MIGRATIONS_PATH -database $DB_URL down
}

function MigrateDrop
{
    migrate -path $MIGRATIONS_PATH -database $DB_URL drop
}

function MigrateForce
{
    if (-not $Version)
    {
        Write-Host 'Usage: .\scripts\make.ps1 migrate-force -Version 3'
        exit 1
    }
    migrate -path $MIGRATIONS_PATH -database $DB_URL force $Version
}

function Lint
{
    golangci-lint run ./...
}

function Fmt
{
    gofmt -s -w .
}

function Vet
{
    go vet ./...
}

function Deps
{
    go mod tidy
    go mod download
}

function Clean
{
    Remove-Item -Recurse -Force $BUILD_DIR -ErrorAction SilentlyContinue
}

function Help
{
    Write-Host 'Usage: .\scripts\make.ps1 [command]'
    Write-Host ''
    Write-Host 'Commands:'
    Write-Host '  build            - собрать проект'
    Write-Host '  run              - запустить проект'
    Write-Host '  dev              - собрать и запустить'
    Write-Host '  test             - запустить тесты'
    Write-Host '  test-verbose     - запустить тесты с выводом'
    Write-Host '  test-cover       - запустить тесты с покрытием'
    Write-Host '  migrate-up       - применить миграции'
    Write-Host '  migrate-down     - откатить миграции'
    Write-Host '  migrate-drop     - дропнуть все миграции'
    Write-Host '  migrate-force    - откатить до версии (-Version N)'
    Write-Host '  lint             - запустить линтер'
    Write-Host '  fmt              - форматировать код'
    Write-Host '  vet              - проверить код'
    Write-Host '  deps             - обновить зависимости'
    Write-Host '  clean            - удалить bin/'
}

# ##################################################
# ROUTER
# ##################################################

switch ($Command)
{
    'build'         {
        Build
    }
    'run'           {
        Run
    }
    'dev'           {
        Dev
    }
    'test'          {
        Test
    }
    'test-verbose'  {
        TestVerbose
    }
    'test-cover'    {
        TestCover
    }
    'migrate-up'    {
        MigrateUp
    }
    'migrate-down'  {
        MigrateDown
    }
    'migrate-drop'  {
        MigrateDrop
    }
    'migrate-force' {
        MigrateForce
    }
    'lint'          {
        Lint
    }
    'fmt'           {
        Fmt
    }
    'vet'           {
        Vet
    }
    'deps'          {
        Deps
    }
    'clean'         {
        Clean
    }
    default         {
        Help
    }
}
