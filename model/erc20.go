package model

import (
	"github.com/ethereum/go-ethereum/common"
)

type ERC20Token struct {
	ContractAddress common.Address
	// Name            string
	Symbol   string
	Decimals uint8
}

func (e *ERC20Token) UpdateSymbol() {
	s := e.Symbol
	str := make([]rune, 0, len(s))
	for _, v := range []rune(s) {
		if v == 0 {
			continue
		}
		str = append(str, v)
	}
	e.Symbol = string(str)
}
