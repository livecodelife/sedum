CREATE TABLE IF NOT EXISTS {{resource|table}} (
    id         SERIAL PRIMARY KEY,
    {{columns}},
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
