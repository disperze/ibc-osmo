package simapp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcgammtypes "github.com/disperze/ibc-osmo/x/intergamm/types"
)

var _ ibcgammtypes.GammKeeper = (*GammKeeperTest)(nil)

type GammKeeperTest struct {
}

func NewGammKeeperTest() *GammKeeperTest {
	return &GammKeeperTest{}
}

func (gamm GammKeeperTest) CalculateSpotPrice(ctx sdk.Context, poolId uint64, tokenInDenom, tokenOutDenom string) (sdk.Dec, error) {
	return sdk.NewDec(1), nil
}

type SwapKeeperTest struct {
}

func NewSwapKeeperTest() *SwapKeeperTest {
	return &SwapKeeperTest{}
}

func (s SwapKeeperTest) SwapExactAmountIn(
	ctx sdk.Context,
	sender sdk.AccAddress,
	poolId uint64,
	tokenIn sdk.Coin,
	tokenOutDenom string,
	tokenOutMinAmount sdk.Int,
) (tokenOutAmount sdk.Int, spotPriceAfter sdk.Dec, err error) {
	return sdk.NewInt(1), sdk.NewDec(1), nil
}
