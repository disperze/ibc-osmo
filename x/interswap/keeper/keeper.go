package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/disperze/ibc-osmo/x/interswap/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey sdk.StoreKey

	transferKeeper types.TransferKeeper

	authKeeper types.AccountKeeper
	gammKeeper types.GammKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	transferKeeper types.TransferKeeper,
	authKeeper types.AccountKeeper,
	gammKeeper types.GammKeeper,
) Keeper {
	return Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		transferKeeper: transferKeeper,
		authKeeper:     authKeeper,
		gammKeeper:     gammKeeper,
	}
}
