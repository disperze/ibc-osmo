package types

// ValidateBasic is used for validating the packet
func (p SpotPricePacketData) ValidateBasic() error {

	// TODO: Validate the packet data

	return nil
}

// GetBytes is a helper for serialising
func (p SpotPricePacketData) GetBytes() ([]byte, error) {
	var modulePacket GammPacketData

	modulePacket.Packet = &GammPacketData_SpotPricePacket{&p}

	return modulePacket.Marshal()
}
