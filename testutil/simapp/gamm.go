package simapp

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	ibcgammtypes "github.com/disperze/ibc-osmo/x/intergamm/types"
)

var _ ibcgammtypes.GammKeeper = (*SwapKeeperTest)(nil)

const InvalidDenom = "n/a"

type SwapKeeperTest struct {
	bankKeeper bankkeeper.Keeper
	// to save swap funds
	moduleName string
}

func NewSwapKeeperTest(bankeeper bankkeeper.Keeper, moduleName string) *SwapKeeperTest {
	return &SwapKeeperTest{bankeeper, moduleName}
}

func (s SwapKeeperTest) SwapExactAmountIn(
	ctx sdk.Context,
	sender sdk.AccAddress,
	poolId uint64,
	tokenIn sdk.Coin,
	tokenOutDenom string,
	tokenOutMinAmount sdk.Int,
) (tokenOutAmount sdk.Int, spotPriceAfter sdk.Dec, err error) {
	if tokenOutDenom == InvalidDenom {
		err = fmt.Errorf("invalid out denom: %s", tokenOutDenom)
		return
	}

	err = s.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, s.moduleName, sdk.NewCoins(tokenIn))
	if err != nil {
		return
	}

	tokensOut := sdk.NewCoins(sdk.NewCoin(tokenOutDenom, tokenIn.Amount))
	err = s.bankKeeper.MintCoins(ctx, s.moduleName, tokensOut)
	if err != nil {
		return
	}

	err = s.bankKeeper.SendCoinsFromModuleToAccount(ctx, s.moduleName, sender, tokensOut)
	if err != nil {
		return
	}

	return tokenIn.Amount, sdk.OneDec(), nil
}
