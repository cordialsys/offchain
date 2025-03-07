package offchain

type GetDepositAddressArgs struct {
	symbol       SymbolId
	network      NetworkId
	accountMaybe AccountName
}

func NewGetDepositAddressArgs(coin SymbolId, network NetworkId, options ...GetDepositAddressOption) GetDepositAddressArgs {
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

func (args *GetDepositAddressArgs) GetSymbol() SymbolId {
	return args.symbol
}

func (args *GetDepositAddressArgs) GetNetwork() NetworkId {
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
