package offchain

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

type Client interface {
	ListAssets() ([]*Asset, error)
	ListBalances(args GetBalanceArgs) ([]*BalanceDetail, error)
	CreateAccountTransfer(args *AccountTransferArgs) (*TransferStatus, error)
	CreateWithdrawal(args *WithdrawalArgs) (*WithdrawalResponse, error)
}
