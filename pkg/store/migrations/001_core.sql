-- +goose Up
CREATE TABLE auth_requests (
        id text PRIMARY KEY,
        application_id text NOT NULL,
        redirect_uri text NOT NULL,
        state text NOT NULL,
        prompt text not null,
        login_hint text not null,
        max_auth_age_seconds integer not null,
        user_id text not null,
        scopes text not null, -- comma separated
        response_type text not null,
        response_mode text not null,
        nonce text not null,
        code_challenge text not null,
        code_challenge_method text not null,
        complete boolean not null,
        created_at INTEGER NOT NULL, -- since epoch
        authenticated_at INTEGER NOT NULL -- since epoch
    );

-- +goose Down
DROP TABLE auth_requests;