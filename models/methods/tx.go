package methods

// The tx method retrieves information on a single transaction, by its
// identifying hash. Expects a response in the form of a TxResponse.
//type TxRequest struct {
//	BaseRequest
//	Transaction string `json:"transaction,omitempty"`
//	Binary      bool   `json:"binary,omitempty"`
//	MinLedger   int64  `json:"min_ledger,omitempty"`
//	MaxLedger   int64  `json:"max_ledger,omitempty"`
//}

// Response expected from a TxRequest.
//type TxResponse struct {
//	BaseResponse
//	Result      TxResponseResult `json:"result,omitempty"`
//	SearchedAll bool             `json:"searched_all,omitempty"`
//}

//type TxResponseResult struct {
//	Transaction
//	Hash        string                     `json:"hash,omitempty"`
//	LedgerIndex int64                      `json:"ledger_index,omitempty"`
//	Meta        models.TransactionMetadata `json:"meta,omitempty"`
//	Validated   bool                       `json:"validated,omitempty"`
//	Date        int64                      `json:"date,omitempty"`
//}
