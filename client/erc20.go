package client

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/fenghaojiang/uniswap-tools-go/constants"
	"github.com/fenghaojiang/uniswap-tools-go/model"
	"github.com/fenghaojiang/uniswap-tools-go/onchain/generated-go/multicall3"
)

func (c *Clients) AggregatedERC20Token(ctx context.Context, address common.Address) (*model.ERC20Token, error) {
	calls := make([]multicall3.Multicall3Call3, 0)

	var err error
	symbolCall := multicall3.Multicall3Call3{
		Target:       address,
		AllowFailure: true,
	}
	symbolCall.CallData, err = c.contractAbis.ERC20.Pack(constants.SymbolMethod)
	if err != nil {
		return nil, err
	}

	calls = append(calls, symbolCall)

	//totalSupplyCall := multicall3.Multicall3Call3{
	//	Target:       address,
	//	AllowFailure: true,
	//}
	//totalSupplyCall.CallData, err = c.contractAbis.ERC20.Pack(constants.TotalSupplyMethod)
	//if err != nil {
	//	return nil, err
	//}
	//calls = append(calls, totalSupplyCall)

	decimalsCall := multicall3.Multicall3Call3{
		Target:       address,
		AllowFailure: true,
	}
	decimalsCall.CallData, err = c.contractAbis.ERC20.Pack(constants.DecimalsMethod)
	if err != nil {
		return nil, err
	}

	calls = append(calls, decimalsCall)

	results, err := c.AggregatedCalls(ctx, calls)
	if err != nil {
		return nil, err
	}
	if len(results) != 2 {
		return nil, fmt.Errorf("failed to match the result, len of result: %d", len(results))
	}

	erc20Token := new(model.ERC20Token)
	erc20Token.ContractAddress = address

	var symbol string
	err = c.contractAbis.ERC20.UnpackIntoInterface(&symbol, constants.SymbolMethod, results[0].ReturnData)
	if err != nil && !strings.Contains(err.Error(), "would go over slice boundary") {
		return nil, err
	}
	if err != nil && strings.Contains(err.Error(), "would go over slice boundary") {
		symbol = strings.Trim(string(results[0].ReturnData), string(rune(0)))
	}
	var decimals uint8
	err = c.contractAbis.ERC20.UnpackIntoInterface(&decimals, constants.DecimalsMethod, results[1].ReturnData)
	if err != nil {
		return nil, err
	}

	return &model.ERC20Token{
		ContractAddress: address,
		Decimals:        decimals,
		Symbol:          symbol,
	}, nil
}

func (c *Clients) ERC20Symbol(ctx context.Context, address common.Address) (string, error) {
	_calldata, err := c.contractAbis.ERC20.Pack(constants.SymbolMethod)
	if err != nil {
		return "", err
	}

	data, err := c.Call(ctx, model.CallContractParam{
		To:       address.String(),
		CallData: hexutil.Encode(_calldata),
	})
	if err != nil {
		return "", err
	}

	var symbol string
	err = c.contractAbis.ERC20.UnpackIntoInterface(&symbol, constants.SymbolMethod, data)
	if err != nil {
		return "", err
	}

	return symbol, nil
}

func (c *Clients) ERC20TotalSupply(ctx context.Context, address common.Address) (*big.Int, error) {
	_calldata, err := c.contractAbis.ERC20.Pack(constants.TotalSupplyMethod)
	if err != nil {
		return nil, err
	}

	data, err := c.Call(ctx, model.CallContractParam{
		To:       address.String(),
		CallData: hexutil.Encode(_calldata),
	})
	if err != nil {
		return nil, err
	}

	totalSupply := new(big.Int)
	err = c.contractAbis.ERC20.UnpackIntoInterface(&totalSupply, constants.TotalSupplyMethod, data)
	if err != nil {
		return nil, err
	}

	return totalSupply, nil
}

func (c *Clients) ERC20Decimals(ctx context.Context, address common.Address) (uint8, error) {
	_calldata, err := c.contractAbis.ERC20.Pack(constants.DecimalsMethod)
	if err != nil {
		return 0, err
	}

	data, err := c.Call(ctx, model.CallContractParam{
		To:       address.String(),
		CallData: hexutil.Encode(_calldata),
	})
	if err != nil {
		return 0, err
	}

	var decimals uint8
	err = c.contractAbis.ERC20.UnpackIntoInterface(&decimals, constants.DecimalsMethod, data)
	if err != nil {
		return 0, err
	}

	return decimals, nil
}

func (c *Clients) ERC20Balance(ctx context.Context, tokenAddress, accountAddress common.Address) (*big.Int, error) {
	_calldata, err := c.contractAbis.ERC20.Pack(constants.BalanceOfMethod, accountAddress)
	if err != nil {
		return nil, err
	}

	data, err := c.Call(ctx, model.CallContractParam{
		To:       tokenAddress.String(),
		CallData: hexutil.Encode(_calldata),
	})
	if err != nil {
		return nil, err
	}

	var balance *big.Int
	err = c.contractAbis.ERC20.UnpackIntoInterface(&balance, constants.BalanceOfMethod, data)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
