package client

import (
	oc "github.com/cordialsys/offchain"
)

type WithdrawalArgs struct {
	address oc.Address
	symbol  oc.SymbolId
	network oc.NetworkId
	amount  oc.Amount
}

func NewWithdrawalArgs(address oc.Address, symbol oc.SymbolId, network oc.NetworkId, amount oc.Amount) WithdrawalArgs {
	return WithdrawalArgs{
		address,
		symbol,
		network,
		amount,
	}
}

func (args *WithdrawalArgs) GetAddress() oc.Address {
	return args.address
}

func (args *WithdrawalArgs) GetSymbol() oc.SymbolId {
	return args.symbol
}

func (args *WithdrawalArgs) GetNetwork() oc.NetworkId {
	return args.network
}

func (args *WithdrawalArgs) GetAmount() oc.Amount {
	return args.amount
}
