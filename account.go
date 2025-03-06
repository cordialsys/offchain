package offchain

type AccountTransferArgs struct {
	from   AccountName
	to     AccountName
	symbol SymbolId
	amount Amount
}

func NewAccountTransferArgs(from AccountName, to AccountName, symbol SymbolId, amount Amount) *AccountTransferArgs {
	return &AccountTransferArgs{
		from:   from,
		to:     to,
		symbol: symbol,
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
