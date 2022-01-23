package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewSpotPricePacketData contructs a new GammPacketData instance
func NewSpotPricePacketData(poolID uint64, tokenIn, tokenOut string) GammPacketData {
	return GammPacketData{
		Packet: &GammPacketData_SpotPrice{
			SpotPrice: &SpotPricePacketData{
				PoolID:   poolID,
				TokenIn:  tokenIn,
				TokenOut: tokenOut,
			},
		},
	}
}

// ValidateBasic is used for validating the packet
func (p SpotPricePacketData) ValidateBasic() error {

	if p.TokenIn == "" {
		return fmt.Errorf("invalid token in denom")
	}

	if p.TokenOut == "" {
		return fmt.Errorf("invalid token out denom")
	}

	return nil
}

// GetBytes is a helper for serialising
func (gpd GammPacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&gpd))
}
