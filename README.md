# CF-Back

API backend **Go** avec **Fiber**, **SQLC** et **Supabase** (base PostgreSQL + Auth).

## Prérequis

- [Go 1.22+](https://go.dev/dl/)
- [SQLC](https://docs.sqlc.dev/en/latest/overview/install.html) (CLI pour générer le code SQL)
- Un projet [Supabase](https://supabase.com) (base + Auth)

### Installer SQLC (Windows)

```powershell
# Avec Go
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
# Puis ajouter $env:GOPATH\bin (ou $env:LOCALAPPDATA\Go\bin) au PATH
```

Ou télécharger le binaire depuis [GitHub Releases](https://github.com/sqlc-dev/sqlc/releases).

## Configuration

1. Copier les variables d’environnement :

   ```powershell
   copy .env.example .env
   ```

2. Renseigner dans `.env` :
   - **SUPABASE_URL** : URL du projet (ex. `https://xxxx.supabase.co`)
   - **SUPABASE_ANON_KEY** : clé « anon public » (Supabase Dashboard > Project Settings > API)
   - **DATABASE_URL** : chaîne de connexion PostgreSQL (Dashboard > Project Settings > Database > Connection string > URI)

3. Appliquer les migrations sur la base Supabase :
   - Fichiers dans `internal/db/migrations/`
   - À exécuter dans l’éditeur SQL du Dashboard Supabase (ou via `supabase db push` si vous utilisez le CLI Supabase)

## Commandes

```powershell
# Télécharger les dépendances
go mod tidy

# Générer le code SQLC (requêtes type-safe)
sqlc generate

# Lancer l’API
go run ./cmd/server
```

L’API écoute sur `http://localhost:3000` (ou la valeur de `PORT`).

## Structure

- **`cmd/server`** : point d’entrée (main)
- **`internal/config`** : chargement de la config (env)
- **`internal/db`** : connexion DB, schéma SQLC, migrations, requêtes SQL
- **`internal/auth`** : middleware Supabase (vérification du JWT via `GET /auth/v1/user`)

## Auth Supabase

Les routes sous `/api/*` exigent un JWT Supabase valide :

- Header : `Authorization: Bearer <access_token>`
- Le token est vérifié en appelant Supabase `GET /auth/v1/user` (recommandé par Supabase)
- L’utilisateur est ensuite disponible dans le contexte (ex. `auth.GetUser(c)`)

Exemples de routes :

- `GET /health` : public
- `GET /api/me` : protégé, renvoie l’utilisateur courant

## SQLC

- **Schéma** : `internal/db/schema/` (tables utilisées par SQLC)
- **Migrations** : `internal/db/migrations/` (à exécuter sur Supabase)
- **Requêtes** : `internal/db/queries/*.sql`
- **Code généré** : `internal/db/sqlc/` (à régénérer avec `sqlc generate` après toute modification du schéma ou des requêtes)
