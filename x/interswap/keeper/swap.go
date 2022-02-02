package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v2/modules/core/04-channel/types"
	"github.com/disperze/ibc-osmo/x/interswap/types"
)

// OnRecvSwapPacket processes packet reception
func (k Keeper) OnRecvSwapPacket(ctx sdk.Context, packet channeltypes.Packet, data types.SwapExactAmountInPacketData) (packetAck types.SwapExactAmountInAck, err error) {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	// logic

	return packetAck, nil
}
