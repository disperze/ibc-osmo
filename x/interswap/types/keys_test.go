package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/disperze/ibc-osmo/x/interswap/types"
)

func TestGetFundAddress(t *testing.T) {
	var (
		port1    = "transfer"
		channel1 = "channel"
		port2    = "transfercha"
		channel2 = "nnel"
	)

	address1 := types.GetFundAddress(port1, channel1)
	address2 := types.GetFundAddress(port2, channel2)
	require.NotEqual(t, address1, address2)
}
