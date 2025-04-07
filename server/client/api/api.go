package api

import (
	oc "github.com/cordialsys/offchain"
)

// NewTransfer creates a new transfer request
func NewTransfer(symbol string, amount oc.Amount) *Transfer {
	return &Transfer{
		Symbol: As(symbol),
		Amount: amount.String(),
	}
}

// WithFromType sets the from account type for a transfer
func (t *Transfer) WithFromType(fromType string) *Transfer {
	t.FromType = As(fromType)
	return t
}

// WithToType sets the to account type for a transfer
func (t *Transfer) WithToType(toType string) *Transfer {
	t.ToType = As(toType)
	return t
}

// WithFrom sets the from account for a transfer
func (t *Transfer) WithFrom(from string) *Transfer {
	t.From = As(from)
	return t
}

// WithTo sets the to account for a transfer
func (t *Transfer) WithTo(to string) *Transfer {
	t.To = As(to)
	return t
}

// NewWithdrawal creates a new withdrawal request
func NewWithdrawal(address string, symbol string, network string, amount oc.Amount) *Withdrawal {
	return &Withdrawal{
		Address: address,
		Symbol:  As(symbol),
		Network: As(network),
		Amount:  amount.String(),
	}
}
