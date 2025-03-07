package offchain

type WithdrawalArgs struct {
	address Address
	symbol  SymbolId
	network NetworkId
	amount  Amount
}

func NewWithdrawalArgs(address Address, symbol SymbolId, network NetworkId, amount Amount) *WithdrawalArgs {
	return &WithdrawalArgs{
		address,
		symbol,
		network,
		amount,
	}
}

func (args *WithdrawalArgs) GetAddress() Address {
	return args.address
}

func (args *WithdrawalArgs) GetSymbol() SymbolId {
	return args.symbol
}

func (args *WithdrawalArgs) GetNetwork() NetworkId {
	return args.network
}

func (args *WithdrawalArgs) GetAmount() Amount {
	return args.amount
}
