package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/disperze/ibc-osmo/x/gamm/types"
	"github.com/tendermint/spm/ibckeeper"
)

type (
	Keeper struct {
		*ibckeeper.Keeper
		cdc      codec.BinaryCodec
		storeKey sdk.StoreKey
		memKey   sdk.StoreKey

		gammKeeper types.GammKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	channelKeeper ibckeeper.ChannelKeeper,
	portKeeper ibckeeper.PortKeeper,
	scopedKeeper ibckeeper.ScopedKeeper,
	gammKeeper types.GammKeeper,
) *Keeper {
	return &Keeper{
		Keeper: ibckeeper.NewKeeper(
			types.PortKey,
			storeKey,
			channelKeeper,
			portKeeper,
			scopedKeeper,
		),
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		gammKeeper: gammKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
