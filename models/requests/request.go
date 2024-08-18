package requests

type BaseRequest struct {
	Id         string `json:"id,omitempty"`
	Command    string `json:"command,omitempty"`
	ApiVersion int16  `json:"api_version,omitempty"`
}

type LookupByLedgerRequest struct {
	LedgerHash  string `json:"ledger_hash,omitempty"`
	LedgerIndex string `json:"ledger_index,omitempty"`
}
