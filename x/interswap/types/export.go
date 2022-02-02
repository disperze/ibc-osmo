package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type InterSwapKeeper interface {
	GetSwapAddress() sdk.AccAddress
}
