package client

import (
	oc "github.com/cordialsys/offchain"
)

type AccountId string

type AccountTransferArgs struct {
	from AccountId
	to   AccountId

	fromType oc.AccountType
	toType   oc.AccountType

	symbol oc.SymbolId
	amount oc.Amount
}

func NewAccountTransferArgs(symbol oc.SymbolId, amount oc.Amount) AccountTransferArgs {
	return AccountTransferArgs{
		"",
		"",
		"",
		"",
		symbol,
		amount,
	}
}

func (args *AccountTransferArgs) SetFrom(from AccountId, fromType oc.AccountType) {
	args.from = from
	args.fromType = fromType
}

func (args *AccountTransferArgs) SetTo(to AccountId, toType oc.AccountType) {
	args.to = to
	args.toType = toType
}

func (args *AccountTransferArgs) GetFrom() (AccountId, oc.AccountType) {
	return args.from, args.fromType
}

func (args *AccountTransferArgs) GetTo() (AccountId, oc.AccountType) {
	return args.to, args.toType
}

func (args *AccountTransferArgs) GetSymbol() oc.SymbolId {
	return args.symbol
}

func (args *AccountTransferArgs) GetAmount() oc.Amount {
	return args.amount
}
