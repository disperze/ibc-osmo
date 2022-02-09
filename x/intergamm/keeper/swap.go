package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v2/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v2/modules/core/04-channel/types"
	"github.com/disperze/ibc-osmo/x/intergamm/types"
)

// OnRecvSwapPacket processes packet reception
func (k Keeper) OnRecvSwapPacket(ctx sdk.Context, packet channeltypes.Packet, sender sdk.AccAddress, amount, denom string, data types.SwapExactAmountInPacketData) (packetAck types.IbcTokenAck, err error) {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	// skip error, validate in transferkeeper Rcv
	tokenIn, _ := types.ParseTokenOnRcv(packet.GetDestPort(), packet.GetDestChannel(), amount, denom)

	// Swap tokens
	tokenOutAmount, tokenOutDenom, err := k.MultihopSwapExactAmountIn(ctx, sender, data.Routes, tokenIn, data.TokenOutMinAmount)
	if err != nil {
		return packetAck, err
	}

	// Send tokens output to source chain
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

func (k Keeper) MultihopSwapExactAmountIn(
	ctx sdk.Context,
	sender sdk.AccAddress,
	routes []types.SwapAmountInRoute,
	tokenIn sdk.Coin,
	tokenOutMinAmount sdk.Int,
) (tokenOutAmount sdk.Int, tokenOutDenom string, err error) {
	for i, route := range routes {
		_outMinAmount := sdk.NewInt(1)
		if len(routes)-1 == i {
			_outMinAmount = tokenOutMinAmount
		}

		tokenOutAmount, _, err = k.gammKeeper.SwapExactAmountIn(ctx, sender, route.PoolId, tokenIn, route.TokenOutDenom, _outMinAmount)
		if err != nil {
			return sdk.Int{}, "", err
		}
		tokenIn = sdk.NewCoin(route.TokenOutDenom, tokenOutAmount)
		tokenOutDenom = route.TokenOutDenom
	}
	return
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
