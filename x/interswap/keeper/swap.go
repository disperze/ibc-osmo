package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v2/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v2/modules/core/04-channel/types"
	"github.com/disperze/ibc-osmo/x/interswap/types"
)

// OnRecvSwapPacket processes packet reception
func (k Keeper) OnRecvSwapPacket(ctx sdk.Context, packet channeltypes.Packet, sender sdk.AccAddress, amount, denom string, data types.SwapExactAmountInPacketData) (packetAck types.SwapExactAmountInAck, err error) {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	// skip error, validate in transferkeeper Rcv
	tokenIn, _ := types.ParseTokenOnRcv(packet.GetDestPort(), packet.GetDestChannel(), amount, denom)
	var tokenOutAmount sdk.Int
	var tokenOutDenom string
	for i, route := range data.Routes {
		_outMinAmount := sdk.NewInt(1)
		if len(data.Routes)-1 == i {
			_outMinAmount = data.TokenOutMinAmount
		}

		tokenOutDenom = route.TokenOutDenom
		tokenOutAmount, _, err = k.gammKeeper.SwapExactAmountIn(ctx, sender, route.PoolId, tokenIn, tokenOutDenom, _outMinAmount)
		if err != nil {
			break
		}
		tokenIn = sdk.NewCoin(tokenOutDenom, tokenOutAmount)
	}

	if err != nil {
		return packetAck, err
	}

	// Send swap output to source chain

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

func (k Keeper) getOutDenomPath(ctx sdk.Context, denom string) (string, error) {
	fullDenomPath := denom

	var err error

	// deconstruct the token denomination into the denomination trace info
	// to determine if the sender is the source chain
	if strings.HasPrefix(denom, "ibc/") {
		fullDenomPath, err = k.transferKeeper.DenomPathFromHash(ctx, denom)
		if err != nil {
			return "", err
		}
	}

	return fullDenomPath, err
}
