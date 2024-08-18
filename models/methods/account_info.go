package methods

import (
	"context"
	"encoding/json"
	"github.com/xrpscan/xrpl-go/models"
	"github.com/xrpscan/xrpl-go/models/ledger"
	"github.com/xrpscan/xrpl-go/models/requests"
)

type AccountInfoRequestOption func(ar *AccountInfoRequest)

func WithLookupByLedger(lr LookupByLedgerRequest) AccountInfoRequestOption {
	// fallback to validated ledger
	if lr.LedgerIndex == "" {
		lr.LedgerIndex = "validated"
	}
	return func(ar *AccountInfoRequest) {
		ar.LookupByLedgerRequest = lr
	}
}

type AccountInfoRequest struct {
	BaseRequest
	LookupByLedgerRequest
	Account     string `json:"account,omitempty"`
	Queue       bool   `json:"queue,omitempty"`
	SignerLists bool   `json:"signer_lists,omitempty"`
	Strict      bool   `json:"strict,omitempty"`
}

func (a *AccountInfoRequest) Context() context.Context {
	return a.BaseRequest.Context()
}

func (a *AccountInfoRequest) Command() string {
	return requests.CommandAccountInfo
}

func (a *AccountInfoRequest) Validate() error {
	//TODO implement validations if needed
	return nil
}

type AccountInfoResponse struct {
	models.BaseResponse
	Result AccountInfoResult `json:"result,omitempty"`
}

type AccountInfoResult struct {
	AccountData ledger.AccountRoot `json:"account_data,omitempty"`

	// A map of account flags parsed out.  This will only be available for rippled nodes 1.11.0 and higher.
	AccountFlags AccountInfoAccountFlags `json:"account_flags,omitempty"`

	// The ledger index of the current in-progress ledger, which was used when
	// retrieving this information.
	LedgerCurrentIndex json.Number `json:"ledger_current_index,omitempty"`

	// The ledger index of the ledger version used when retrieving this
	// information. The information does not contain any changes from ledger
	// versions newer than this one.
	LedgerIndex json.Number `json:"ledger_index,omitempty"`

	// Information about queued transactions sent by this account. This
	// information describes the state of the local rippled server, which may be
	// different from other servers in the peer-to-peer XRP Ledger network. Some
	// fields may be omitted because the values are calculated "lazily" by the
	// queuing mechanism.
	QueueData AccountQueueData `json:"queue_data,omitempty"`

	// True if this data is from a validated ledger version; if omitted or set
	// to false, this data is not final.
	Validated bool `json:"validated,omitempty"`
}

type AccountQueueTransaction struct {

	// * Whether this transaction changes this address's ways of authorizing
	// * transactions.
	AuthChange    bool   `json:"auth_change,omitempty"`
	Fee           string `json:"fee,omitempty"`
	FeeLevel      string `json:"fee_level,omitempty"`
	MaxSpendDrops string `json:"max_spend_drops,omitempty"`
	Seq           int    `json:"seq,omitempty"`
}

type AccountQueueData struct {
	//Number of queued transactions from this address.
	TxnCount int `json:"txn_count,omitempty"`

	// Whether a transaction in the queue changes this address's ways of
	// authorizing transactions. If true, this address can queue no further
	// transactions until that transaction has been executed or dropped from the
	// queue.
	AuthChangeQueued bool `json:"auth_change_queued,omitempty"`

	//The lowest Sequence Number among transactions queued by this address.
	LowestSequence int `json:"lowest_sequence,omitempty"`

	// The highest Sequence Number among transactions queued by this address.
	HighestSequence int `json:"highest_sequence,omitempty"`

	// Integer amount of drops of XRP that could be debited from this address if
	// every transaction in the queue consumes the maximum amount of XRP possible.
	MaxSpendDropsTotal string `json:"max_spend_drops_total,omitempty"`

	// Information about each queued transaction from this address.
	Transactions []AccountQueueTransaction `json:"transactions,omitempty"`
}

type AccountInfoAccountFlags struct {

	// Enable rippling on this address's trust lines by default. Required for issuing addresses; discouraged for others.
	DefaultRipple bool `json:"default_ripple,omitempty"`

	// This account can only receive funds from transactions it sends, and from preauthorized accounts.
	// (It has DepositAuth enabled.)
	DepositAuth bool `json:"deposit_auth,omitempty"`

	// Disallows use of the master key to sign transactions for this account.
	DisableMasterKey bool `json:"disable_master_key,omitempty"`

	// Disallow incoming Checks from other accounts.
	DisallowIncomingCheck bool `json:"disallow_incoming_check,omitempty"`

	// Disallow incoming NFTOffers from other accounts. Part of the DisallowIncoming amendment.
	DisallowIncomingNFTokenOffer bool `json:"disallow_incoming_nf_token_offer,omitempty"`

	// Disallow incoming PayChannels from other accounts. Part of the DisallowIncoming amendment.
	DisallowIncomingPayChan bool `json:"disallow_incoming_pay_chan,omitempty"`

	// Disallow incoming Trustlines from other accounts. Part of the DisallowIncoming amendment.
	DisallowIncomingTrustline bool `json:"disallow_incoming_trustline,omitempty"`

	// Client applications should not send XRP to this account. Not enforced by rippled.
	DisallowIncomingXRP bool `json:"disallow_incoming_xrp,omitempty"`

	// All assets issued by this address are frozen.
	GlobalFreeze bool `json:"global_freeze,omitempty"`

	// This address cannot freeze trust lines connected to it. Once enabled, cannot be disabled.
	NoFreeze bool `json:"no_freeze,omitempty"`

	// The account has used its free SetRegularKey transaction.
	PasswordSpent bool `json:"password_spent,omitempty"`

	// This account must individually approve other users for those users to hold this account's issued currencies.
	RequireAuthorization bool `json:"require_authorization,omitempty"`

	// Requires incoming payments to specify a Destination Tag.
	RequireDestinationTag bool `json:"require_destination_tag,omitempty"`

	// This address can claw back issued IOUs. Once enabled, cannot be disabled.
	AllowTrustLineClawback bool `json:"allow_trust_line_clawback,omitempty"`
}
