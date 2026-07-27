package scim

import (
	"context"
	"time"

	"github.com/idpzero/idpzero/pkg/configuration"
)

// Action describes what happened to a single user during a sync.
const (
	ActionCreated = "created"
	ActionUpdated = "updated"
	ActionFailed  = "failed"
)

// Result is the outcome of provisioning a single user.
type Result struct {
	Subject  string
	UserName string
	Action   string // created | updated | failed
	Err      string // populated when Action == failed
}

// Report is the outcome of a full sync run.
type Report struct {
	RanAt    time.Time
	Endpoint string
	Results  []Result
	Failed   int
}

// syncClient is the subset of *Client used by Sync, extracted so tests can
// substitute a fake.
type syncClient interface {
	FindByExternalID(ctx context.Context, externalID string) (*User, error)
	Create(ctx context.Context, user User) (*User, error)
	Replace(ctx context.Context, id string, user User) (*User, error)
}

// Sync reconciles the given users against the target: each user is looked up by
// externalId and either created (when absent) or replaced (when present). A
// failure on one user is recorded and does not abort the run. RanAt is left
// zero for the caller to stamp.
func Sync(ctx context.Context, c syncClient, endpoint string, users []*configuration.User) Report {
	report := Report{
		Endpoint: endpoint,
		Results:  make([]Result, 0, len(users)),
	}

	for _, u := range users {
		if u == nil {
			continue
		}

		payload := ToSCIMUser(u)
		res := Result{Subject: u.Subject, UserName: payload.UserName}

		existing, err := c.FindByExternalID(ctx, u.Subject)
		if err != nil {
			res.Action = ActionFailed
			res.Err = err.Error()
			report.Failed++
			report.Results = append(report.Results, res)
			continue
		}

		if existing == nil {
			_, err = c.Create(ctx, payload)
			res.Action = ActionCreated
		} else {
			_, err = c.Replace(ctx, existing.ID, payload)
			res.Action = ActionUpdated
		}

		if err != nil {
			res.Action = ActionFailed
			res.Err = err.Error()
			report.Failed++
		}

		report.Results = append(report.Results, res)
	}

	return report
}
