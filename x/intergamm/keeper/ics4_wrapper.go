package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v2/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v2/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v2/modules/core/exported"

	"github.com/disperze/ibc-osmo/x/intergamm/types"
)

type SwapICS4Wrapper struct {
	channelKeeper transfertypes.ChannelKeeper
}

func NewSwapICS4Wrapper(channelKeeper transfertypes.ChannelKeeper) SwapICS4Wrapper {
	return SwapICS4Wrapper{
		channelKeeper: channelKeeper,
	}
}

func (w SwapICS4Wrapper) GetChannel(ctx sdk.Context, srcPort, srcChan string) (channel channeltypes.Channel, found bool) {
	return w.channelKeeper.GetChannel(ctx, srcPort, srcChan)
}

func (w SwapICS4Wrapper) GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool) {
	return w.channelKeeper.GetNextSequenceSend(ctx, portID, channelID)
}

func (w SwapICS4Wrapper) SendPacket(ctx sdk.Context, channelCap *capabilitytypes.Capability, packet ibcexported.PacketI) error {
	swapAddress := types.GetFundAddress(packet.GetSourcePort(), packet.GetSourceChannel()).String()
	var data transfertypes.FungibleTokenPacketData
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal sent packet data: %s", err.Error())
	}

	// if Sent from the swap address, skip packet
	if swapAddress == data.Sender {
		return nil
	}

	return w.channelKeeper.SendPacket(ctx, channelCap, packet)
}
