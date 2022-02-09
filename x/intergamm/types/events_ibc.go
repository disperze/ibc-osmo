package types

// IBC events
const (
	EventTypeSwapPacket     = "swap_packet"
	EventTypeJoinPoolPacket = "join_pool_packet"
	EventTypeExitPoolPacket = "exit_pool_packet"

	AttributeKeyAckSuccess = "success"
	AttributeKeyAck        = "acknowledgement"
	AttributeKeyAckError   = "error"
)
