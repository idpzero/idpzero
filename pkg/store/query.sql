-- name: GetAuthRequestByID :one
SELECT * FROM
  auth_requests
WHERE
  id = ? LIMIT 1;

-- name: GetAuthRequestByAuthCode :one
SELECT * FROM
  auth_requests
WHERE
  auth_code = ? LIMIT 1;

-- name: UpdateAuthRequestUser :execrows
UPDATE auth_requests 
SET user_id = ?, complete = 1, authenticated_at = ? 
WHERE id = ?;

-- name: UpdateAuthCode :execrows
UPDATE auth_requests 
SET auth_code = ?
WHERE id = ?;

-- name: CreateAuthRequest :one
INSERT INTO auth_requests (
    id,
    application_id,
    redirect_uri,
    state,
    prompt,
    login_hint,
    max_auth_age_seconds,
    user_id,
    scopes,
    response_type,
    response_mode,
    code_challenge,
    code_challenge_method,
    nonce,
    complete,
    created_at,
    authenticated_at
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?,?, ?) RETURNING *;

-- name: DeleteAllAuthRequests :exec
DELETE FROM auth_requests;

-- name: CreateToken :one
INSERT INTO tokens (
    id,
    auth_request_id,
    application_id,
    refresh_token_id,
    subject,
    audience,
    expiration,
    scopes,
    created_at
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?,?, ?) RETURNING *;

-- name: GetTokenByID :one
SELECT * FROM
  tokens
WHERE
  id = ? LIMIT 1;


-- name: DeleteAllTokens :exec
DELETE FROM tokens; 
  
-- name: CreateKey :one
INSERT INTO keys (
    id,
    alg,
    usage,
    public_key,
    private_key,
    created_at
  )
VALUES
  (?, ?, ?, ?, ?, ?) RETURNING *;
  
-- name: GetKeyByID :one
SELECT * FROM
  keys
WHERE
  id = ? LIMIT 1;

-- name: GetKeysByUse :many
SELECT * FROM
  keys
WHERE
  usage = ?;

-- name: GetAllKeys :many
SELECT * FROM
  keys;

-- name: DeleteAllKeys :exec
DELETE FROM keys;



-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    id,
    auth_time,
    amr,
    audience,
    subject,
    application_id,
    expiration,
    scopes,
    created_at
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?,?, ?) RETURNING *;


-- name: GetRefreshTokenByID :one
SELECT * FROM
  refresh_tokens
WHERE
  id = ? LIMIT 1;


-- CREATE TABLE refresh_tokens (
--         id text PRIMARY KEY,
--         auth_time INTEGER NOT NULL,
--         amr text NOT NULL,
--         audience text NOT NULL,
--         user_id text NOT NULL,
--         application_id text NOT NULL,
--         expiration INTEGER NOT NULL, -- since epoch
--         scopes text NOT NULL, -- comma separated
--         created_at INTEGER NOT NULL -- since epoch
--     );