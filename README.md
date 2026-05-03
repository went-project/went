# WENT CLI

WENT, Go ve Cobra ile yazilmis bir proje scaffold aracidir. Yeni bir backend proje iskeleti olusturur ve route/model/controller/migration dosyalarini template ile uretir.

> **For AI Agents:** See [CLAUDE.md](./CLAUDE.md) for agent guidelines, which points to [AGENTS.md](./AGENTS.md) for detailed conventions.

## Ozellikler

- Yeni proje olusturma (`create [name]`)
- Hazir uygulama iskeleti (Chi router, GORM, graceful shutdown)
- Swagger dokumantasyonu otomatik olarak olusturulur (`swag init` create sirasinda calisir)
- Controller metodlarinda Swagger annotation desteği (Payload/Response DTO tipleri)
- Varsayilan `User` route/model/controller uretilmesi
- `up` ve `down` SQL migration dosyasi uretilmesi
- SQLite, MySQL ve PostgreSQL destegi
- `.env` dosyasi yoksa varsayilan DB davranisi SQLite fallback olarak calisir
- Ayrik generator komutlari:
  - `gen:router [name]`
  - `gen:model [name]`
  - `gen:controller [name]`
  - `gen:migration [name]`
  - `gen:resource [name]`
  - Tum `gen:*` komutlarinda varsayilan davranis: iliskili dosyalari skeleton olarak uretmek
  - Tum `gen:*` komutlarinda `-a/--all`: iliskili dosyalari full (dolu) template ile uretmek
- Dahili migration runner (harici araç gerekmez):
  - `migrate` — bekleyen migration'ları uygular
  - `migrate:rollback` — son N migration'ı geri alır
  - `migrate:fresh` — tüm tabloları düşürüp sıfırdan uygular

## Gereksinimler

- Go 1.24+
- [swag](https://github.com/swaggo/swag) CLI (`go install github.com/swaggo/swag/cmd/swag@latest`)

## Kurulum ve Build

```bash
# Projeyi derle
make build

# Binary kur (varsayilan: /usr/local/bin/went2)
make install
```

Alternatif:

```bash
go build -o build/went2 .
```

## Komutlar

```bash
# Yardim
./build/went2 --help

# Versiyon
./build/went2 version

# Yeni proje olustur
./build/went2 create my-app

# Generator komutlari (proje dizini icerisinde calistirilmali)
./build/went2 gen:model User             # skeleton related set
./build/went2 gen:model User -a          # full related set
./build/went2 gen:controller User        # skeleton related set
./build/went2 gen:controller User -a     # full related set
./build/went2 gen:router User            # skeleton related set
./build/went2 gen:router User -a         # full related set
./build/went2 gen:resource User          # skeleton related set
./build/went2 gen:resource User -a       # full related set
./build/went2 gen:migration User         # skeleton related set
./build/went2 gen:migration User -a      # full related set

# Migration runner komutlari (proje dizini icerisinde calistirilmali)
./build/went2 migrate                     # bekleyen tum up dosyalarini siraya gore uygular
./build/went2 migrate:rollback            # son 1 migration'i geri alir
./build/went2 migrate:rollback --step 3   # son 3 migration'i geri alir
./build/went2 migrate:fresh               # tum down'lari calistirip sifirdan up uygular
```

## Release Pipeline

GitHub Actions release workflow'u `.github/workflows/release.yml` icerisinde tanimlidir. `main` branch'ine push geldiginde pipeline otomatik olarak yeni bir `v1.0526.<number>` tag'i uretir ve onu repoya push eder. Tag push edildiginde ise release buildleri calisir:

- Windows: `went2-windows-amd64-v.1.0526.1.exe`
- Linux: `went2-linux-amd64-v.1.0526.1`, `went2-linux-arm64-v.1.0526.1`
- macOS: `went2-darwin-amd64-v.1.0526.1`, `went2-darwin-arm64-v.1.0526.1`

Ornek tag:

```bash
git tag v1.0526.1
git push origin v1.0526.1
```

`main` push akisi, mevcut ay-yil prefiksine gore bir sonraki numarayi hesaplar ve `v1.0526.<number>` formatinda tag olusturur. Tag push oldugunda ayni workflow release artifact'larini uretir ve GitHub Release olarak yayinlar.

## Kurulum

Release'ten indirmek icin kendi platformuna uygun scripti calistir:

```bash
curl -fsSL https://raw.githubusercontent.com/went-project/went/main/install.sh -o install.sh
chmod +x install.sh
./install.sh
```

Windows icin:

```powershell
Invoke-WebRequest https://raw.githubusercontent.com/went-project/went/main/install.ps1 -OutFile install.ps1
./install.ps1
```

Varsayilan olarak en son release indirilir. Belirli bir tag kurmak istersen `WENT_VERSION=v1.0526.1 ./install.sh` ya da PowerShell tarafinda `$env:WENT_VERSION='v1.0526.1'; ./install.ps1` kullanabilirsin.

Kurulum sonunda binary kullanici PATH'ine kalici olarak eklenir ve terminalde `went2 --help` ile dogrulayabilirsin.

> **On Kosul:** `create` disindaki tum komutlar calisma dizininde `wentconfig.json` dosyasi bekler. Dosya yoksa komut isleme devam etmez.

### migrate:rollback Flag'leri

| Flag | Kisaltma | Varsayilan | Aciklama |
|---|---|---|---|
| `--step` | `-s` | `1` | Kac adim geri alinacagini belirtir |

### gen:model Flag'leri

| Flag | Kisaltma | Aciklama |
|---|---|---|
| `--migration` | `-m` | Geriye donuk uyumluluk icin korunmustur (skeleton related mode zaten migration olusturur) |
| `--all` | `-a` | Tum iliskili dosyalari full template icerigiyle uretir |

## Generation Modlari

- Varsayilan (`gen:* <Name>`): model, migration, request, controller, router, resource dosyalarini skeleton (bos/iskelet) olarak uretir.
- Full (`gen:* <Name> --all`): ayni iliskili dosyalari full template icerigiyle uretir.

## Create Ile Uretilenler

`create my-app` komutu asagidakileri olusturur:

- Proje klasor yapisi (`internal`, `http`, `routes`, `database`, vb.)
- `main.go` — Chi server, graceful shutdown, Swagger docs import
- `internal/config/config.go` — env tabanli konfigurasyon (SQLite/MySQL/Postgres)
- `internal/providers/database_provider.go` — GORM baglantisi
- `internal/responses/error_response.go` — global JSON error response helper
- `http/middlewares/auth.go`
- `http/middlewares/cors.go`
- `routes/main_router.go`
- `routes/user_router.go`
- `database/models/user.go`
- `http/requests/user_request.go` — create/update request DTO'lari
- `http/controllers/user_controller.go` — Swagger annotation'li CRUD, helper metodlar
- `http/resources/pagination.go` — ortak pagination meta yapisi
- `http/resources/user_resource.go` — resource DTO + query builder (paginated/tekil okuma)
- `database/migrations/000001_create_users_table.up.sql`
- `database/migrations/000001_create_users_table.down.sql`
- `.env.example`, `.gitignore`, `.devwatchignore`
- `docs/` — `swag init` ile otomatik olusturulur

Create sirasinda su komutlar otomatik calisir:
1. `go mod init <name>`
2. `go mod tidy`
3. `swag init`

## Ortam Degiskenleri

Olusturulan proje `.env.local` dosyasini kullanir (`APP_ENV=local`).

```env
APP_ENV=local
PORT=8080

# SQLite icin:
DB_DIALECT=sqlite
DB_NAME=./database.sqlite

# PostgreSQL icin:
# DB_DIALECT=postgres
# DB_HOST=localhost
# DB_PORT=5432
# DB_USER=postgres
# DB_PASSWORD=password
# DB_NAME=mydb

JWT_SECRET=changeme
```

> **Not:** SQLite kullaniminda `DB_NAME` alanini mutlaka belirtin. `DB_STORAGE` yazilmissa `DB_NAME` olarak otomatik kullanilir.

## Controller Yapisi

Uretilen controller'lar asagidaki yapiya sahiptir:

- Ortak `ParseID(r)` helper'i `http/controllers/helpers.go` icinde bulunur ve tum controller'lar tarafindan kullanilir
- `findByID(id, &item)` — DB lookup helper
- Controller icinde local payload/response struct tanimlamaz
- Create/Update payload tiplerini `http/requests` altindan kullanir (`requests.XPayload`, `requests.XUpdatePayload`)
- Response tiplerini `http/resources` altindan kullanir (`resources.XResource`, `resources.XCollection`)
- Her metod icin tam Swagger annotation (`@Summary`, `@Tags`, `@Param`, `@Success`, `@Router`)
- `GetAll` ve `GetByID` okuma sorgularini `http/resources` altindaki `XQuery` yapisina delege eder
- Hata donuslerinde `internal/responses` altindaki global helper fonksiyonlarini kullanir (`JSONError`, `JSONErrorWithDetails`)

## Resource Yapisi

Uretilen resource dosyalari (`http/resources/*_resource.go`) asagidaki yapiya sahiptir:

- `XResource` — tek kayit response DTO'su
- `XCollection` — Laravel Collection benzeri response envelope:
  - `data`: donusen kayit listesi
  - `meta`: pagination bilgisi (`current_page`, `last_page`, `per_page`, `total`, `from`, `to`)
- `XQuery` — resource icinde DB okuma sorgularini yoneten katman:
  - `Paginate(page, perPage)`
  - `Find(id)`

## Global Error Response Yapisi

`create` ile uretilen `internal/responses/error_response.go` ortak hata payload'ini saglar:

- `ErrorBody` — standart error json body (`success`, `message`, `details`)
- `JSONError(...)` — standart hata donusu
- `JSONErrorWithDetails(...)` — detayli hata donusu

## Migration Runner Detaylari

`migrate`, `migrate:rollback` ve `migrate:fresh` komutlari:

- Calistirilacak dizindeki `.env.local` veya `.env` dosyasini okur (DB_DIALECT, DB_NAME, vb.)
- `wentmigrations` adli bir takip tablosu olusturur ve yönetir
- `database/migrations/` altindaki `*.up.sql` ve `*.down.sql` dosyalarini dosya adi siralamasina gore calistirir
- SQLite, MySQL ve PostgreSQL ile calisir

## Notlar

- Generator ve migrate komutlari mevcut proje klasoru icinde calistirilmalidir.
- `gen:router <Name>` cikti dosyasini `routes/<name>_router.go` formatinda (kucuk harf) uretir.
- `gen:migration` komutu `.up.sql` ve `.down.sql` dosyasi uretir.
- Router/controller template'leri app adini `wentconfig.json` uzerinden okur.
- `go get -u` create akisindan cikarilmistir; Go surum uyumsuzluklarini onlemek icin.
- Hata mesajlari artik altta yatan hata detayini da icerir.
