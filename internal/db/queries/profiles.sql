-- name: GetProfileByID :one
SELECT id, email, display_name, avatar_url, created_at, updated_at
FROM public.profiles
WHERE id = $1;

-- name: CreateProfile :one
INSERT INTO public.profiles (id, email, display_name, avatar_url)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateProfile :one
UPDATE public.profiles
SET
    display_name = COALESCE(sqlc.narg('display_name'), display_name),
    avatar_url = COALESCE(sqlc.narg('avatar_url'), avatar_url),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;
