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
        authenticated_at INTEGER NOT NULL, -- since epoch
        auth_code text 
    );

CREATE TABLE tokens (
        id text PRIMARY KEY,
        auth_request_id text, -- originator if available
        application_id text NOT NULL,
        refresh_token_id text,
        subject text NOT NULL,
        audience text NOT NULL,
        expiration INTEGER NOT NULL, -- since epoch
        scopes text NOT NULL, -- comma separated
        created_at INTEGER NOT NULL -- since epoch
    );

CREATE TABLE refresh_tokens (
        id text PRIMARY KEY,
        auth_time INTEGER NOT NULL,
        amr text,
        audience text NOT NULL,
        subject text NOT NULL,
        application_id text NOT NULL,
        expiration INTEGER NOT NULL, -- since epoch
        scopes text NOT NULL, -- comma separated
        created_at INTEGER NOT NULL -- since epoch
    );

CREATE TABLE keys (
        id text PRIMARY KEY,
        alg text NOT NULL,
        usage text NOT NULL,
        public_key text NOT NULL,
        private_key text NOT NULL,
        created_at INTEGER NOT NULL -- since epoch
    );

-- +goose Down
DROP TABLE auth_requests;
DROP TABLE tokens;
DROP TABLE keys;