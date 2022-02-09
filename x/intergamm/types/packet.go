package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewIbcPacketData contructs a new IbcPacketData instance
func NewIbcSwapPacketData(sender, receiver, amount, denom string, routes []SwapAmountInRoute, tokenOut sdk.Int) IbcPacketData {
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

// NewIbcJoinPoolPacketData contructs a new IbcPacketData instance with JoinPool msg
func NewIbcJoinPoolPacketData(sender, receiver, amount, denom string, poolID uint64, shareOutMinAmount sdk.Int) IbcPacketData {
	return IbcPacketData{
		Amount:   amount,
		Denom:    denom,
		Sender:   sender,
		Receiver: receiver,
		Gamm: &IbcPacketData_JoinPool{
			JoinPool: &JoinPoolPacketData{
				Sender:            sender,
				PoolId:            poolID,
				ShareOutMinAmount: shareOutMinAmount,
			},
		},
	}
}

// NewIbcJoinPoolPacketData contructs a new IbcPacketData instance with JoinPool msg
func NewIbcExitPoolPacketData(sender, receiver, amount, denom string, tokenOutDenom string, tokenOutMinAmount sdk.Int) IbcPacketData {
	return IbcPacketData{
		Amount:   amount,
		Denom:    denom,
		Sender:   sender,
		Receiver: receiver,
		Gamm: &IbcPacketData_ExitPool{
			ExitPool: &ExitPoolPacketData{
				Sender:            sender,
				TokenOutDenom:     tokenOutDenom,
				TokenOutMinAmount: tokenOutMinAmount,
			},
		},
	}
}

// GetBytes is a helper for serialising
func (gpd IbcPacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&gpd))
}
