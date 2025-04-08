package api

import (
	"fmt"
	"strings"

	oc "github.com/cordialsys/offchain"
)

type SymbolAndChain string

// OKX reports the chain with the symbol included, e.g. "BTC-Bitcoin", or "ETH-Ethereum"
func (s SymbolAndChain) Split() (oc.SymbolId, oc.NetworkId) {
	parts := strings.Split(string(s), "-")
	symbolId := parts[0]
	chainId := symbolId
	if len(parts) > 1 {
		chainId = parts[1]
	}
	return oc.SymbolId(symbolId), oc.NetworkId(chainId)
}

func (s SymbolAndChain) SymbolId() oc.SymbolId {
	symbolId, _ := s.Split()
	return symbolId
}

func (s SymbolAndChain) NetworkId() oc.NetworkId {
	_, networkId := s.Split()
	return networkId
}

func NewSymbolAndChain(symbol oc.SymbolId, network oc.NetworkId) SymbolAndChain {
	return SymbolAndChain(fmt.Sprintf("%s-%s", symbol, network))
}
