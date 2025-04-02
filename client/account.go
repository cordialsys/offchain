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

var CoreFunding AccountName = "core/funding"
var CoreTrading AccountName = "core/trading"
var CoreMargin AccountName = "core/margin"

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
	from   AccountName
	to     AccountName
	symbol oc.SymbolId
	amount oc.Amount
}

func NewAccountTransferArgs(from AccountName, to AccountName, symbol oc.SymbolId, amount oc.Amount) AccountTransferArgs {
	return AccountTransferArgs{
		from,
		to,
		symbol,
		amount,
	}
}

func (args *AccountTransferArgs) GetFrom() AccountName {
	return args.from
}

func (args *AccountTransferArgs) GetTo() AccountName {
	return args.to
}

func (args *AccountTransferArgs) GetSymbol() oc.SymbolId {
	return args.symbol
}

func (args *AccountTransferArgs) GetAmount() oc.Amount {
	return args.amount
}
