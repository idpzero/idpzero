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


-- name: DeleteAllKeys :exec
DELETE FROM keys;
