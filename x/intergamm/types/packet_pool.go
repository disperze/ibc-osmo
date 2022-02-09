package types

import (
	"fmt"
)

// ValidateBasic is used for validating the packet
func (p JoinPoolPacketData) ValidateBasic() error {

	if p.Sender == "" {
		return fmt.Errorf("invalid address")
	}

	if !p.ShareOutMinAmount.IsPositive() {
		return fmt.Errorf("invalid share out min amount")
	}

	return nil
}

// ValidateBasic is used for validating the packet
func (p ExitPoolPacketData) ValidateBasic() error {

	if p.Sender == "" {
		return fmt.Errorf("invalid address")
	}

	if p.TokenOutDenom == "" {
		return fmt.Errorf("invalid token out denom")
	}

	if !p.TokenOutMinAmount.IsPositive() {
		return fmt.Errorf("invalid out min amount")
	}

	return nil
}
