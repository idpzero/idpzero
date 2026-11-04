-- +goose Up
-- acr and amr of the signed-in user, copied from their configured auth on login.
ALTER TABLE auth_requests ADD COLUMN acr text NOT NULL DEFAULT '';
ALTER TABLE auth_requests ADD COLUMN amr text NOT NULL DEFAULT ''; -- space separated

-- +goose Down
ALTER TABLE auth_requests DROP COLUMN amr;
ALTER TABLE auth_requests DROP COLUMN acr;
