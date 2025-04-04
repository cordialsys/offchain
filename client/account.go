package client

import (
	oc "github.com/cordialsys/offchain"
)

type AccountTransferArgs struct {
	from oc.AccountId
	to   oc.AccountId

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

func (args *AccountTransferArgs) SetFrom(from oc.AccountId) {
	args.from = from
}

func (args *AccountTransferArgs) SetTo(to oc.AccountId) {
	args.to = to
}

func (args *AccountTransferArgs) SetFromType(fromType oc.AccountType) {
	args.fromType = fromType
}

func (args *AccountTransferArgs) SetToType(toType oc.AccountType) {
	args.toType = toType
}

func (args *AccountTransferArgs) GetTo() (oc.AccountId, oc.AccountType) {
	return args.to, args.toType
}

func (args *AccountTransferArgs) GetFrom() (oc.AccountId, oc.AccountType) {
	return args.from, args.fromType
}

func (args *AccountTransferArgs) GetSymbol() oc.SymbolId {
	return args.symbol
}

func (args *AccountTransferArgs) GetAmount() oc.Amount {
	return args.amount
}

// transferring between the same subaccount or main account?
func (args *AccountTransferArgs) IsSameAccount() bool {
	return args.from == args.to
}
