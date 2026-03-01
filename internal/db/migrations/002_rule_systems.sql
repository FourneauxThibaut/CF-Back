-- Rule systems: Phase 1 JSON storage (rules + block_definitions in JSONB)
CREATE TABLE IF NOT EXISTS public.rule_systems (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    rules JSONB NOT NULL DEFAULT '[]',
    block_definitions JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE public.rule_systems ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can manage own rule systems"
    ON public.rule_systems FOR ALL
    USING (auth.uid() = user_id)
    WITH CHECK (auth.uid() = user_id);
