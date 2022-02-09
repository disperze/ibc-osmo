package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewIbcPacketData contructs a new IbcPacketData instance
func NewIbcPacketData(sender, receiver, amount, denom string, routes []SwapAmountInRoute, tokenOut sdk.Int) IbcPacketData {
	return IbcPacketData{
		Amount:   amount,
		Denom:    denom,
		Sender:   sender,
		Receiver: receiver,
		Gamm: &IbcPacketData_Swap{
			Swap: &SwapExactAmountInPacketData{
				Sender:            sender,
				Routes:            routes,
				TokenOutMinAmount: tokenOut,
			},
		},
	}
}

// GetBytes is a helper for serialising
func (gpd IbcPacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&gpd))
}

// GetBytes is a helper for serialising
func (gpd IbcPacketData) GetSafeBytes() ([]byte, error) {
	bz, err := ModuleCdc.MarshalJSON(&gpd)
	if err != nil {
		return nil, err
	}

	return sdk.SortJSON(bz)
}
