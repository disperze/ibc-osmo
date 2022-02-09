package types

import (
	"fmt"
)

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
