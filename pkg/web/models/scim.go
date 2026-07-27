package models

import (
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/idpzero/idpzero/pkg/scim"
)

// SCIMResultRow is a single displayable per-user provisioning outcome.
type SCIMResultRow struct {
	Subject  string
	UserName string
	Action   string // created | updated | failed
	Error    string
}

// SCIMModel is the view model for the provisioning page.
type SCIMModel struct {
	Configured      bool // an endpoint is set
	Endpoint        string
	TokenConfigured bool
	UserCount       int

	HasRun  bool
	RanAt   string
	Total   int
	Failed  int
	Results []SCIMResultRow
}

// NewSCIMModel builds the provisioning view model from the SCIM config and the
// last sync report (which may be nil when a sync has never run).
func NewSCIMModel(cfg *configuration.SCIMConfig, tokenConfigured bool, userCount int, report *scim.Report) SCIMModel {
	model := SCIMModel{
		TokenConfigured: tokenConfigured,
		UserCount:       userCount,
	}

	if cfg != nil && cfg.Endpoint != "" {
		model.Configured = true
		model.Endpoint = cfg.Endpoint
	}

	if report == nil {
		return model
	}

	model.HasRun = true
	model.RanAt = report.RanAt.Format("2006-01-02 15:04:05 MST")
	model.Total = len(report.Results)
	model.Failed = report.Failed
	model.Results = make([]SCIMResultRow, 0, len(report.Results))
	for _, res := range report.Results {
		model.Results = append(model.Results, SCIMResultRow{
			Subject:  res.Subject,
			UserName: res.UserName,
			Action:   res.Action,
			Error:    res.Err,
		})
	}

	return model
}
