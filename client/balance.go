package client

type GetBalanceArgs struct {
	accountType AccountType
}

func NewGetBalanceArgs(accountType AccountType) GetBalanceArgs {
	return GetBalanceArgs{
		accountType: accountType,
	}
}

func (args *GetBalanceArgs) GetAccountType() AccountType {
	return args.accountType
}
