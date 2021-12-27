package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GammKeeper defines the expected osmosis gamm keeper
type GammKeeper interface {
	CalculateSpotPrice(ctx sdk.Context, poolId uint64, tokenInDenom, tokenOutDenom string) (sdk.Dec, error)
}
