package offchain

type TransactionId string
type OperationStatus string

const (
	OperationStatusPending OperationStatus = "pending"
	OperationStatusSuccess OperationStatus = "success"
	OperationStatusFailed  OperationStatus = "failed"
)

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
	Symbol        SymbolId          `json:"symbol"`
	Network       NetworkId         `json:"network"`
	Amount        Amount            `json:"amount"`
	Fee           Amount            `json:"fee"`
	TransactionId TransactionId     `json:"transaction_id,omitempty"`
	Comment       string            `json:"comment,omitempty"`
	Notes         map[string]string `json:"notes,omitempty"`
}

type Client interface {
	ListAssets() ([]*Asset, error)
	ListBalances(args GetBalanceArgs) ([]*BalanceDetail, error)
	CreateAccountTransfer(args AccountTransferArgs) (*TransferStatus, error)
	CreateWithdrawal(args WithdrawalArgs) (*WithdrawalResponse, error)
	GetDepositAddress(args GetDepositAddressArgs) (Address, error)
	ListWithdrawalHistory(args WithdrawalHistoryArgs) ([]*WithdrawalHistory, error)
}
