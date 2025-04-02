package client

type GetBalanceArgs struct {
	account AccountName
}

func NewGetBalanceArgs(account AccountName) GetBalanceArgs {
	return GetBalanceArgs{
		account: account,
	}
}

func (args *GetBalanceArgs) GetAccount() AccountName {
	return args.account
}
