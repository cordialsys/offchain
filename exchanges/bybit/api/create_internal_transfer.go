package api

import (
	oc "github.com/cordialsys/offchain"
	"github.com/google/uuid"
)

type TransferRequest struct {
	TransferID      string         `json:"transferId"`      // UUID, must be manually generated
	Coin            oc.SymbolId    `json:"coin"`            // Coin symbol in uppercase
	Amount          string         `json:"amount"`          // Transfer amount
	FromAccountType oc.AccountType `json:"fromAccountType"` // Source account type
	ToAccountType   oc.AccountType `json:"toAccountType"`   // Destination account type
}

type TransferStatus string

const TransferStateSuccess TransferStatus = "SUCCESS"

type TransferResult struct {
	TransferID string         `json:"transferId"`
	Status     TransferStatus `json:"status"`
}
type TransferResponse = Response[TransferResult]

// https://bybit-exchange.github.io/docs/v5/asset/transfer/create-inter-transfer
func (c *Client) CreateInternalTransfer(coin oc.SymbolId, amount oc.Amount, fromAccount oc.AccountType, toAccount oc.AccountType) (*TransferResponse, error) {
	uid := uuid.New().String()
	request := TransferRequest{
		TransferID:      uid,
		Coin:            coin,
		Amount:          amount.String(),
		FromAccountType: fromAccount,
		ToAccountType:   toAccount,
	}
	var response TransferResponse
	_, err := c.Request("POST", "/v5/asset/transfer/inter-transfer", &request, &response, nil)
	if err != nil {
		return nil, err
	}
	return &response, nil

}
