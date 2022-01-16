package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v2/modules/core/04-channel/types"
	"github.com/disperze/ibc-osmo/x/intergamm/types"
)

// OnRecvSpotPricePacket processes packet reception
func (k Keeper) OnRecvSpotPricePacket(ctx sdk.Context, packet channeltypes.Packet, data types.SpotPricePacketData) (packetAck types.SpotPricePacketAck, err error) {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	price, err := k.gammKeeper.CalculateSpotPrice(ctx, data.PoolId, data.TokenIn, data.TokenOut)
	if err != nil {
		return packetAck, err
	}

	packetAck.Price = price.String()

	return packetAck, nil
}
