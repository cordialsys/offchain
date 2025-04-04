package api

import (
	oc "github.com/cordialsys/offchain"
	"github.com/google/uuid"
)

type UniversalTransferRequest struct {
	TransferID      string         `json:"transferId"`      // UUID, must be manually generated
	Coin            oc.SymbolId    `json:"coin"`            // Coin symbol in uppercase
	Amount          string         `json:"amount"`          // Transfer amount
	FromMemberId    int            `json:"fromMemberId"`    // From UID
	ToMemberId      int            `json:"toMemberId"`      // To UID
	FromAccountType oc.AccountType `json:"fromAccountType"` // Source account type
	ToAccountType   oc.AccountType `json:"toAccountType"`   // Destination account type
}

// https://bybit-exchange.github.io/docs/v5/asset/transfer/unitransfer
func (c *Client) CreateUniversalTransfer(coin oc.SymbolId, amount oc.Amount, fromMemberId int, toMemberId int, fromAccountType oc.AccountType, toAccountType oc.AccountType) (*TransferResponse, error) {
	uid := uuid.New().String()
	request := UniversalTransferRequest{
		TransferID:      uid,
		Coin:            coin,
		Amount:          amount.String(),
		FromMemberId:    fromMemberId,
		ToMemberId:      toMemberId,
		FromAccountType: fromAccountType,
		ToAccountType:   toAccountType,
	}
	var response TransferResponse
	_, err := c.Request("POST", "/v5/asset/transfer/universal-transfer", &request, &response, nil)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
