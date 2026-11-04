-- +goose Up
-- acr and amr of the signed-in user, copied from their configured auth on login.
ALTER TABLE auth_requests ADD COLUMN acr text NOT NULL DEFAULT '';
ALTER TABLE auth_requests ADD COLUMN amr text NOT NULL DEFAULT ''; -- space separated
-- acr of the login a refresh token continues (amr is already kept).
ALTER TABLE refresh_tokens ADD COLUMN acr text NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE refresh_tokens DROP COLUMN acr;
ALTER TABLE auth_requests DROP COLUMN amr;
ALTER TABLE auth_requests DROP COLUMN acr;
