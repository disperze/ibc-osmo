package types

import "fmt"

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
