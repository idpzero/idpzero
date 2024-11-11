// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package query

import (
	"context"
)

const createAuthRequest = `-- name: CreateAuthRequest :one
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
  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?,?, ?) RETURNING id, application_id, redirect_uri, state, prompt, login_hint, max_auth_age_seconds, user_id, scopes, response_type, response_mode, nonce, code_challenge, code_challenge_method, complete, created_at, authenticated_at
`

type CreateAuthRequestParams struct {
	ID                  string
	ApplicationID       string
	RedirectUri         string
	State               string
	Prompt              string
	LoginHint           string
	MaxAuthAgeSeconds   int64
	UserID              string
	Scopes              string
	ResponseType        string
	ResponseMode        string
	CodeChallenge       string
	CodeChallengeMethod string
	Nonce               string
	Complete            bool
	CreatedAt           int64
	AuthenticatedAt     int64
}

func (q *Queries) CreateAuthRequest(ctx context.Context, arg CreateAuthRequestParams) (AuthRequest, error) {
	row := q.db.QueryRowContext(ctx, createAuthRequest,
		arg.ID,
		arg.ApplicationID,
		arg.RedirectUri,
		arg.State,
		arg.Prompt,
		arg.LoginHint,
		arg.MaxAuthAgeSeconds,
		arg.UserID,
		arg.Scopes,
		arg.ResponseType,
		arg.ResponseMode,
		arg.CodeChallenge,
		arg.CodeChallengeMethod,
		arg.Nonce,
		arg.Complete,
		arg.CreatedAt,
		arg.AuthenticatedAt,
	)
	var i AuthRequest
	err := row.Scan(
		&i.ID,
		&i.ApplicationID,
		&i.RedirectUri,
		&i.State,
		&i.Prompt,
		&i.LoginHint,
		&i.MaxAuthAgeSeconds,
		&i.UserID,
		&i.Scopes,
		&i.ResponseType,
		&i.ResponseMode,
		&i.Nonce,
		&i.CodeChallenge,
		&i.CodeChallengeMethod,
		&i.Complete,
		&i.CreatedAt,
		&i.AuthenticatedAt,
	)
	return i, err
}

const getAuthRequestByID = `-- name: GetAuthRequestByID :one
SELECT id, application_id, redirect_uri, state, prompt, login_hint, max_auth_age_seconds, user_id, scopes, response_type, response_mode, nonce, code_challenge, code_challenge_method, complete, created_at, authenticated_at FROM
  auth_requests
WHERE
  id = ? LIMIT 1
`

func (q *Queries) GetAuthRequestByID(ctx context.Context, id string) (AuthRequest, error) {
	row := q.db.QueryRowContext(ctx, getAuthRequestByID, id)
	var i AuthRequest
	err := row.Scan(
		&i.ID,
		&i.ApplicationID,
		&i.RedirectUri,
		&i.State,
		&i.Prompt,
		&i.LoginHint,
		&i.MaxAuthAgeSeconds,
		&i.UserID,
		&i.Scopes,
		&i.ResponseType,
		&i.ResponseMode,
		&i.Nonce,
		&i.CodeChallenge,
		&i.CodeChallengeMethod,
		&i.Complete,
		&i.CreatedAt,
		&i.AuthenticatedAt,
	)
	return i, err
}
