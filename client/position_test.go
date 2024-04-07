package client

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"

	"github.com/fenghaojiang/uniswap-tools-go/constants"
)

func TestOnAccountHoldings(t *testing.T) {
	clis, err := NewClientsWithEndpoints([]string{
		"https://rpc.ankr.com/polygon",
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	position, err := clis.WithNetwork(constants.PolygonNetwork).AggregatedPosition(ctx, []*big.Int{
		new(big.Int).SetInt64(869899),
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v", position)
}

func TestOnAccountHoldingsEthereum(t *testing.T) {
	clis, err := NewClientsWithEndpoints([]string{
		"https://rpc.ankr.com/eth",
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	results, err := clis.WithNetwork(constants.EthereumNetwork).AggregatedPosition(ctx, []*big.Int{
		new(big.Int).SetInt64(324342),
	})
	assert.NoError(t, err)
	for _, item := range results {
		t.Log("lockedAmount0", item.LockedAmount0.String())
		t.Log("lockedAmount1", item.LockedAmount1.String())
		t.Log("rewards0amount", item.FeeRewards0Amount.String())
		t.Log("rewards1amount", item.FeeRewards1Amount.String())
	}
}
