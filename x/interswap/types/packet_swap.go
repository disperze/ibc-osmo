package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewSwapGammPacket contructs a new GammPacketData instance
func NewSwapGammPacket(sender string, routes []SwapAmountInRoute, tokenOut sdk.Int) GammPacketData {
	return GammPacketData{
		Packet: &GammPacketData_Swap{
			Swap: &SwapExactAmountInPacketData{
				Sender:            sender,
				Routes:            routes,
				TokenOutMinAmount: tokenOut,
			},
		},
	}
}

// ValidateBasic is used for validating the packet
func (p SwapExactAmountInPacketData) ValidateBasic() error {

	if p.Sender == "" {
		return fmt.Errorf("invalid address")
	}

	if len(p.Routes) == 0 {
		return fmt.Errorf("invalid routes")
	}

	if !p.TokenOutMinAmount.IsPositive() {
		return fmt.Errorf("invalid token out min amount")
	}

	return nil
}

// GetBytes is a helper for serialising
func (gpd GammPacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&gpd))
}

// GetBytes is a helper for serialising
func (gpd GammPacketData) GetSafeBytes() ([]byte, error) {
	bz, err := ModuleCdc.MarshalJSON(&gpd)
	if err != nil {
		return nil, err
	}

	return sdk.SortJSON(bz)
}
