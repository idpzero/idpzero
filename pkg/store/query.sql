-- name: GetAuthRequestByID :one
SELECT * FROM
  auth_requests
WHERE
  id = ? LIMIT 1;

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