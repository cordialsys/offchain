package offchain

import (
	"path/filepath"
	"strings"
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
	symbol SymbolId
	amount Amount
}

func NewAccountTransferArgs(from AccountName, to AccountName, symbol SymbolId, amount Amount) AccountTransferArgs {
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

func (args *AccountTransferArgs) GetSymbol() SymbolId {
	return args.symbol
}

func (args *AccountTransferArgs) GetAmount() Amount {
	return args.amount
}
