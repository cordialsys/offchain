package client

import (
	oc "github.com/cordialsys/offchain"
)

type TransactionId string
type OperationStatus string

const (
	OperationStatusPending OperationStatus = "pending"
	OperationStatusSuccess OperationStatus = "success"
	OperationStatusFailed  OperationStatus = "failed"
)

type BalanceDetail struct {
	// These are the original exchange IDs
	SymbolId oc.SymbolId `json:"symbol_id"`
	// often the network is not relevant, as assets can be withdrawn to multiple networks
	NetworkId oc.NetworkId `json:"network_id,omitempty"`

	// Amount, accounted for decimals
	Available   oc.Amount `json:"available"`
	Unavailable oc.Amount `json:"unavailable"`
}

type WithdrawalResponse struct {
	// any ID returned by the exchange
	ID     string
	Status OperationStatus
}

type TransferStatus struct {
	ID     string
	Status OperationStatus
}

type WithdrawalHistory struct {
	ID            string            `json:"id"`
	Status        OperationStatus   `json:"status"`
	Symbol        oc.SymbolId       `json:"symbol"`
	Network       oc.NetworkId      `json:"network"`
	Amount        oc.Amount         `json:"amount"`
	Fee           oc.Amount         `json:"fee"`
	TransactionId TransactionId     `json:"transaction_id,omitempty"`
	Comment       string            `json:"comment,omitempty"`
	Notes         map[string]string `json:"notes,omitempty"`
}

type Client interface {
	// List all of the support assets on the exchange
	ListAssets() ([]*oc.Asset, error)

	// List of all of the balances on an account
	ListBalances(args GetBalanceArgs) ([]*BalanceDetail, error)

	// Create a transfer between accounts (e.g. main to sub, sub to sub, etc) on the exchange
	CreateAccountTransfer(args AccountTransferArgs) (*TransferStatus, error)

	// Withdraw funds to an external wallet
	CreateWithdrawal(args WithdrawalArgs) (*WithdrawalResponse, error)

	// Get a deposit address for an asset
	GetDepositAddress(args GetDepositAddressArgs) (oc.Address, error)

	// List paginated withdrawal history on an account in descending order
	ListWithdrawalHistory(args WithdrawalHistoryArgs) ([]*WithdrawalHistory, error)
}
