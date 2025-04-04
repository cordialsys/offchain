package client

import oc "github.com/cordialsys/offchain"

type GetBalanceArgs struct {
	accountType oc.AccountType
}

func NewGetBalanceArgs(accountType oc.AccountType) GetBalanceArgs {
	return GetBalanceArgs{
		accountType: accountType,
	}
}

func (args *GetBalanceArgs) GetAccountType() oc.AccountType {
	return args.accountType
}

func (args *GetBalanceArgs) SetAccountType(accountType oc.AccountType) {
	args.accountType = accountType
}
