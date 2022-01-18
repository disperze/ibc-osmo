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
