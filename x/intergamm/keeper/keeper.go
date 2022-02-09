package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/disperze/ibc-osmo/x/intergamm/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey sdk.StoreKey

	transferKeeper types.TransferKeeper
	gammKeeper     types.GammKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	transferKeeper types.TransferKeeper,
	gammKeeper types.GammKeeper,
) Keeper {
	return Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		transferKeeper: transferKeeper,
		gammKeeper:     gammKeeper,
	}
}
