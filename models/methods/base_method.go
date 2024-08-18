package methods

import (
	"context"
)

type LookupByLedgerRequest struct {
	LedgerHash  string `json:"ledger_hash,omitempty"`
	LedgerIndex string `json:"ledger_index,omitempty"`
}

type BaseRequest struct {
	// ctx is unexported to prevent people from using Context wrong
	// and mutating the contexts held by callers of the same request.
	ctx        context.Context
	Id         string `json:"id,omitempty"`
	Command    string `json:"command,omitempty"`
	ApiVersion int16  `json:"api_version,omitempty"`
}

func (a *BaseRequest) Context() context.Context {
	if a.ctx != nil {
		return a.ctx
	}
	return context.Background()
}

type BaseResponse struct {
	Id         string            `json:"id,omitempty"`
	Status     string            `json:"status,omitempty"`
	Type       string            `json:"type,omitempty"`
	Warning    string            `json:"warning,omitempty"`
	Warnings   []ResponseWarning `json:"warnings,omitempty"`
	Forwarded  bool              `json:"forwarded,omitempty"`
	ApiVersion int16             `json:"api_version,omitempty"`
}

type ResponseWarning struct {
	Id      int               `json:"id,omitempty"`
	Message string            `json:"message,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}
