package ledger

import "encoding/json"

type BaseLedgerEntry struct {
	index string
}

type HasPreviousTxnID struct {

	/**
	 * The identifying hash of the transaction that most recently modified this
	 * object.
	 */
	PreviousTxnID string
	/**
	 * The index of the ledger that contains the transaction that most recently
	 * modified this object.
	 */
	PreviousTxnLgrSeq json.Number
}

type HasOptionalPreviousTxnID struct {

	/**
	 * The identifying hash of the transaction that most recently modified this
	 * object. This field was added in the `fixPreviousTxnID` amendment, so it
	 * may not be present in every object.
	 */
	PreviousTxnID string
	/**
	 * The index of the ledger that contains the transaction that most recently
	 * modified this object. This field was added in the `fixPreviousTxnID`
	 * amendment, so it may not be present in every object.
	 */
	PreviousTxnLgrSeq json.Number
}
