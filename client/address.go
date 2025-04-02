package client

import (
	oc "github.com/cordialsys/offchain"
)

type GetDepositAddressArgs struct {
	symbol       oc.SymbolId
	network      oc.NetworkId
	accountMaybe AccountName
}

func NewGetDepositAddressArgs(coin oc.SymbolId, network oc.NetworkId, options ...GetDepositAddressOption) GetDepositAddressArgs {
	args := GetDepositAddressArgs{
		coin,
		network,
		"",
	}
	for _, option := range options {
		option(&args)
	}
	return args
}

func (args *GetDepositAddressArgs) GetSymbol() oc.SymbolId {
	return args.symbol
}

func (args *GetDepositAddressArgs) GetNetwork() oc.NetworkId {
	return args.network
}

func (args *GetDepositAddressArgs) GetAccount() (AccountName, bool) {
	return args.accountMaybe, args.accountMaybe != ""
}

func (args *GetDepositAddressArgs) GetAccountId() (string, bool) {
	id := args.accountMaybe.Id()
	return id, id != ""
}

type GetDepositAddressOption func(args *GetDepositAddressArgs)

func WithAccount(account AccountName) GetDepositAddressOption {
	return func(args *GetDepositAddressArgs) {
		args.accountMaybe = account
	}
}
