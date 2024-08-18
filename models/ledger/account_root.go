package ledger

import "encoding/json"

type AccountRoot struct {
	BaseLedgerEntry
	HasPreviousTxnID
	LedgerEntryType string

	// The identifying (classic) address of this account.
	Account string

	// The account's current XRP balance in drops, represented as a string.
	Balance string

	// A bit-map of boolean flags enabled for this account.
	Flags json.Number

	//The number of objects this account owns in the ledger, which contributes
	//to its owner reserve.
	OwnerCount json.Number

	// The sequence number of the next valid transaction for this account.
	Sequence json.Number

	//The identifying hash of the transaction most recently sent by this
	//account. This field must be enabled to use the AccountTxnID transaction
	//field. To enable it, send an AccountSet transaction with the.
	//`asfAccountTxnID` flag enabled
	AccountTxnID string

	//The ledger entry ID of the corresponding AMM ledger entry.
	//Set during account creation; cannot be modified.
	//If present, indicates that this is a special AMM AccountRoot; always omitted on non-AMM accounts.
	AMMID string

	//A domain associated with this account. In JSON, this is the hexadecimal
	//for the ASCII representation of the domain.
	Domain string

	// The md5 hash of an email address.
	EmailHash string

	//A public key that may be used to send encrypted messages to this account
	//in JSON, uses hexadecimal.
	MessageKey string

	//The address of a key pair that can be used to sign transactions for this
	//account instead of the master key. Use a SetRegularKey transaction to
	//change this value.
	RegularKey string

	//How many Tickets this account owns in the ledger. This is updated
	//automatically to ensure that the account stays within the hard limit of 250.
	//Tickets at a time.
	TicketCount json.Number

	//How many significant digits to use for exchange rates of Offers involving
	//currencies issued by this address. Valid values are 3 to 15, inclusive.
	TickSize json.Number

	//A transfer fee to charge other users for sending currency issued by this
	//account to each other.
	TransferRate json.Number

	// An arbitrary 256-bit value that users can set.
	WalletLocator string

	// Total NFTokens this account's issued that have been burned. This number is always equal or less than MintedNFTokens.
	BurnedNFTokens json.Number

	// The sequence that the account first minted an NFToken
	FirstNFTSequence json.Number

	// Total NFTokens have been minted by and on behalf of this account.
	MintedNFTokens json.Number

	// Another account that can mint NFTokens on behalf of this account.
	NFTokenMinter string
}
