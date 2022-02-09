package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v2/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v2/modules/core/04-channel/types"
	"github.com/disperze/ibc-osmo/x/intergamm/types"
)

// OnRecvJoinPoolPacket processes packet reception
func (k Keeper) OnRecvJoinPoolPacket(ctx sdk.Context, packet channeltypes.Packet, sender sdk.AccAddress, amount, denom string, data types.JoinPoolPacketData) (packetAck types.JoinPoolAck, err error) {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	// skip error, validate in transferkeeper Rcv
	tokenIn, _ := types.ParseTokenOnRcv(packet.GetDestPort(), packet.GetDestChannel(), amount, denom)

	// Swap tokens
	tokenOutAmount, err := k.gammKeeper.JoinSwapExternAmountIn(ctx, sender, data.PoolId, tokenIn, data.ShareOutMinAmount)
	if err != nil {
		return packetAck, err
	}

	// Send tokens output to source chain
	tokenOutDenom := types.GetPoolShareDenom(data.PoolId)
	tokenTransferOut := sdk.NewCoin(tokenOutDenom, tokenOutAmount)
	// transferKeeper needs ICS4Wrapper
	err = k.transferKeeper.SendTransfer(ctx, packet.GetDestPort(), packet.GetDestChannel(), tokenTransferOut, sender, data.Sender, clienttypes.ZeroHeight(), 0)
	if err != nil {
		return packetAck, err
	}

	denomPathOut, err := k.getOutDenomPath(ctx, tokenOutDenom)
	if err != nil {
		return packetAck, err
	}

	packetAck.Amount = tokenOutAmount.String()
	packetAck.Denom = denomPathOut

	return packetAck, nil
}
