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
	return strings.TrimPrefix(string(name), "accounts/")
}

func NewAccountName(id string) AccountName {
	return AccountName(filepath.Join("accounts", id))
}

type GetBalanceArgs struct {
	account AccountName
}

func NewGetBalanceArgs(account AccountName) *GetBalanceArgs {
	return &GetBalanceArgs{
		account: account,
	}
}

func (args *GetBalanceArgs) GetAccount() AccountName {
	return args.account
}
