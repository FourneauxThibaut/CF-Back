-- Schema initial pour CF-Back
-- Les utilisateurs sont gérés par Supabase Auth (auth.users)
-- Exemple: table profil liée à l'id Supabase

-- Extensions si besoin (Supabase les fournit souvent déjà)
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Exemple: profil utilisateur (optionnel, à adapter)
CREATE TABLE IF NOT EXISTS public.profiles (
    id UUID PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
    email TEXT,
    display_name TEXT,
    avatar_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- RLS (Row Level Security) recommandé avec Supabase
ALTER TABLE public.profiles ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can read own profile"
    ON public.profiles FOR SELECT
    USING (auth.uid() = id);

CREATE POLICY "Users can update own profile"
    ON public.profiles FOR UPDATE
    USING (auth.uid() = id);

CREATE POLICY "Users can insert own profile"
    ON public.profiles FOR INSERT
    WITH CHECK (auth.uid() = id);
