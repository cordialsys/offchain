package offchain

type WithdrawalHistoryArgs struct {
	limit         int
	nextPageToken string
}

func NewWithdrawalHistoryArgs() WithdrawalHistoryArgs {
	return WithdrawalHistoryArgs{
		limit:         100,
		nextPageToken: "",
	}
}
func (args WithdrawalHistoryArgs) WithPageToken(nextPageToken string) WithdrawalHistoryArgs {
	args.nextPageToken = nextPageToken
	return args
}

func (args WithdrawalHistoryArgs) WithLimit(limit int) WithdrawalHistoryArgs {
	args.limit = limit
	return args
}

func (args *WithdrawalHistoryArgs) GetLimit() int {
	return args.limit
}

func (args *WithdrawalHistoryArgs) GetPageToken() string {
	return args.nextPageToken
}

func (args *WithdrawalHistoryArgs) SetPageToken(nextPageToken string) {
	args.nextPageToken = nextPageToken
}

func (args *WithdrawalHistoryArgs) SetLimit(limit int) {
	args.limit = limit
}
