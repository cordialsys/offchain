package client

import (
	"path/filepath"
	"strings"

	oc "github.com/cordialsys/offchain"
)

// - accounts/funding
// - accounts/trading
// - accounts/x/subacounts/y

// - core/funding
// - core/trading
// - core/margin

type AccountName string
type AccountType string

// var CoreFunding AccountName = "core/funding"
// var CoreTrading AccountName = "core/trading"
// var CoreMargin AccountName = "core/margin"

func (name AccountName) Id() string {
	if strings.HasPrefix(string(name), "core/") {
		return ""
	}
	return strings.TrimPrefix(string(name), "accounts/")
}

func NewAccountName(id string) AccountName {
	return AccountName(filepath.Join("accounts", id))
}

type AccountTransferArgs struct {
	from AccountName
	to   AccountName

	fromType AccountType
	toType   AccountType

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

func (args *AccountTransferArgs) SetFrom(from AccountName, fromType AccountType) {
	args.from = from
	args.fromType = fromType
}

func (args *AccountTransferArgs) SetTo(to AccountName, toType AccountType) {
	args.to = to
	args.toType = toType
}

func (args *AccountTransferArgs) GetFrom() (AccountName, AccountType) {
	return args.from, args.fromType
}

func (args *AccountTransferArgs) GetTo() (AccountName, AccountType) {
	return args.to, args.toType
}

func (args *AccountTransferArgs) GetSymbol() oc.SymbolId {
	return args.symbol
}

func (args *AccountTransferArgs) GetAmount() oc.Amount {
	return args.amount
}
