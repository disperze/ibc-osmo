package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	transfertypes "github.com/cosmos/ibc-go/v2/modules/apps/transfer/types"
)

func ParseDenonOnRcv(sourcePort, sourceChannel, denom string) string {
	if transfertypes.ReceiverChainIsSource(sourcePort, sourceChannel, denom) {
		voucherPrefix := transfertypes.GetDenomPrefix(sourcePort, sourceChannel)
		unprefixedDenom := denom[len(voucherPrefix):]

		denomTrace := transfertypes.ParseDenomTrace(unprefixedDenom)
		if denomTrace.Path != "" {
			return denomTrace.IBCDenom()
		}

		return unprefixedDenom
	}

	sourcePrefix := transfertypes.GetDenomPrefix(sourcePort, sourceChannel)
	// NOTE: sourcePrefix contains the trailing "/"
	prefixedDenom := sourcePrefix + denom

	denomTrace := transfertypes.ParseDenomTrace(prefixedDenom)

	return denomTrace.IBCDenom()
}

func ParseTokenOnRcv(sourcePort, sourceChannel, amount, denom string) (sdk.Coin, error) {
	fullDenom := ParseDenonOnRcv(sourcePort, sourceChannel, denom)
	transferAmount, ok := sdk.NewIntFromString(amount)
	if !ok {
		return sdk.Coin{}, sdkerrors.Wrapf(transfertypes.ErrInvalidAmount, "unable to parse transfer amount (%s) into sdk.Int", amount)
	}

	return sdk.NewCoin(fullDenom, transferAmount), nil
}
