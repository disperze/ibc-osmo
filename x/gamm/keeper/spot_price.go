package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	"github.com/disperze/ibc-osmo/x/gamm/types"
)

// OnRecvSpotPricePacket processes packet reception
func (k Keeper) OnRecvSpotPricePacket(ctx sdk.Context, packet channeltypes.Packet, data types.SpotPricePacketData) (packetAck types.SpotPricePacketAck, err error) {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	// TODO: packet reception logic

	return packetAck, nil
}

// OnAcknowledgementSpotPricePacket responds to the the success or failure of a packet
// acknowledgement written on the receiving chain.
func (k Keeper) OnAcknowledgementSpotPricePacket(ctx sdk.Context, packet channeltypes.Packet, data types.SpotPricePacketData, ack channeltypes.Acknowledgement) error {
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:

		// TODO: failed acknowledgement logic
		_ = dispatchedAck.Error

		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck types.SpotPricePacketAck

		if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		// TODO: successful acknowledgement logic

		return nil
	default:
		// The counter-party module doesn't implement the correct acknowledgment format
		return errors.New("invalid acknowledgment format")
	}
}

// OnTimeoutSpotPricePacket responds to the case where a packet has not been transmitted because of a timeout
func (k Keeper) OnTimeoutSpotPricePacket(ctx sdk.Context, packet channeltypes.Packet, data types.SpotPricePacketData) error {

	// TODO: packet timeout logic

	return nil
}
