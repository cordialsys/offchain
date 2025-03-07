package offchain

type Client interface {
	ListAssets() ([]*Asset, error)
	ListBalances(args GetBalanceArgs) ([]*BalanceDetail, error)
	CreateAccountTransfer(args *AccountTransferArgs) error
	CreateWithdrawal(args *WithdrawalArgs) error
}
