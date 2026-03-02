# CF-Back

API Go (Gin + SQLC + Supabase). PostgreSQL pour la DB, Supabase pour l’auth.

## Setup

Go 1.22+, [SQLC](https://docs.sqlc.dev/en/latest/overview/install.html). Sur Windows : `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest` puis ajouter le bin au PATH.

```powershell
copy .env.example .env
# Éditer .env avec SUPABASE_URL, SUPABASE_ANON_KEY, DATABASE_URL
```

Migrations dans `internal/db/migrations/` — à exécuter dans le SQL Editor Supabase (ou `supabase db push`).

```powershell
go mod tidy
sqlc generate
go run .
```

Écoute sur `:3000` (ou `PORT`).

## Auth

Routes `/api/*` protégées par JWT Supabase. Header : `Authorization: Bearer <token>`. Le middleware appelle `GET /auth/v1/user` et expose l’user via `auth.GetUser(c)`.

- `GET /health` — public
- `GET /api/me`, `GET /api/profile` — protégés

## Scalingo

Même variables que le dev : `DATABASE_URL`, `SUPABASE_URL`, `SUPABASE_ANON_KEY`. Les définir dans Dashboard > Environment. Pour CORS, ajouter `FRONTEND_URL` si le front est sur un autre domaine.

```bash
scalingo --app cherry-fire env-set DATABASE_URL="..." SUPABASE_URL="..." SUPABASE_ANON_KEY="..."
```

Utiliser l’URI pooler (port 6543) pour la connexion DB — plus adapté au serverless.

## Structure

- `main.go`, `router/`, `handlers/` — bootstrap et routes
- `internal/config` — config depuis env
- `internal/db` — pool pgx, schéma SQLC, migrations
- `internal/auth` — middleware JWT Supabase
- `internal/db/sqlc/` — code généré, à régénérer avec `sqlc generate` après modif du schéma ou des queries
