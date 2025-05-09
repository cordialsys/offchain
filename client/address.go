package client

import (
	oc "github.com/cordialsys/offchain"
)

type GetDepositAddressArgs struct {
	symbol     oc.SymbolId
	network    oc.NetworkId
	subaccount oc.AccountId
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

func (args *GetDepositAddressArgs) GetSubaccount() (oc.AccountId, bool) {
	return args.subaccount, args.subaccount != ""
}

type GetDepositAddressOption func(args *GetDepositAddressArgs)

func WithSubaccount(subaccount oc.AccountId) GetDepositAddressOption {
	return func(args *GetDepositAddressArgs) {
		args.subaccount = subaccount
	}
}

func (args *GetDepositAddressArgs) SetSubaccount(subaccount oc.AccountId) {
	args.subaccount = subaccount
}
