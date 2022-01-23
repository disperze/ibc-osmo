package intergamm_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v2/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v2/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v2/modules/core/24-host"
	ibctesting "github.com/cosmos/ibc-go/v2/testing"

	"github.com/disperze/ibc-osmo/testutil/simapp"
	"github.com/disperze/ibc-osmo/x/intergamm/types"
)

type GammTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
}

func (suite *GammTestSuite) SetupTest() {
	ibctesting.DefaultTestingAppInit = simapp.SetupTestingApp

	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(1))
}

func NewInterGammPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = types.PortID
	path.EndpointB.ChannelConfig.PortID = types.PortID
	path.EndpointA.ChannelConfig.Version = types.Version
	path.EndpointB.ChannelConfig.Version = types.Version

	return path
}

// constructs a send from chainA to chainB on the established channel/connection
// and sends the same coin back from chainB to chainA.
func (suite *GammTestSuite) TestOnRecvPacket() {
	var packetData []byte
	testCases := []struct {
		name          string
		malleate      func()
		expAckSuccess bool
	}{
		{
			"success", func() {}, true,
		},
		{
			"ack - invalid denom a", func() {
				spotPricePacket := types.NewSpotPricePacketData("ID321", 1, "", "uion")
				packetData = spotPricePacket.GetBytes()
			}, false,
		},
		{
			"ack - invalid denom b", func() {
				spotPricePacket := types.NewSpotPricePacketData("", 1, "uosmo", "")
				packetData = spotPricePacket.GetBytes()
			}, false,
		},
		{
			"ack - cannot unmarshal packet data", func() {
				packetData = []byte("invalid data")
			}, false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			// setup between chainA and chainB
			path := NewInterGammPath(suite.chainA, suite.chainB)

			suite.coordinator.Setup(path)
			timeoutHeight := clienttypes.NewHeight(0, 100)
			seq := uint64(1)

			spotPricePacket := types.NewSpotPricePacketData("juno1q4aw0vtcr4jdj70g", 1, "uosmo", "uion")
			packetData = spotPricePacket.GetBytes()

			tc.malleate()

			packet := channeltypes.NewPacket(packetData, seq, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
			err := packet.ValidateBasic()
			suite.Require().NoError(err)

			err = path.EndpointA.SendPacket(packet)
			suite.Require().NoError(err)

			ackBytes, err := suite.RecvPacket(path.EndpointB, packet)
			suite.Require().NoError(err)

			ctx := suite.chainB.GetContext()
			_, found := path.EndpointB.Chain.App.GetIBCKeeper().ChannelKeeper.GetPacketAcknowledgement(ctx, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, seq)
			suite.Require().True(found)

			jsonAck := string(ackBytes)
			suite.Require().Equal(tc.expAckSuccess, !strings.Contains(jsonAck, "error"))
		})
	}
}

// RecvPacket define custom endpoint.RecvPacket to get Acknowledgement
func (suite *GammTestSuite) RecvPacket(endpoint *ibctesting.Endpoint, packet channeltypes.Packet) ([]byte, error) {
	// get proof of packet commitment on source
	packetKey := host.PacketCommitmentKey(packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())
	proof, proofHeight := endpoint.Counterparty.Chain.QueryProof(packetKey)

	recvMsg := channeltypes.NewMsgRecvPacket(packet, proof, proofHeight, endpoint.Chain.SenderAccount.GetAddress().String())

	// receive on counterparty and update source client
	r, err := endpoint.Chain.SendMsgs(recvMsg)
	if err != nil {
		return nil, err
	}

	err = endpoint.Counterparty.UpdateClient()
	if err != nil {
		return nil, err
	}

	ackBytes, exists := FindAcknowledgement(r)
	if !exists {
		return nil, fmt.Errorf("could not find ack")
	}

	return ackBytes, nil
}

func FindAcknowledgement(txResult *sdk.Result) ([]byte, bool) {
	for _, event := range txResult.Events {
		if event.Type != channeltypes.EventTypeWriteAck {
			continue
		}

		for _, attr := range event.Attributes {
			if string(attr.Key) == channeltypes.AttributeKeyAck {
				return attr.Value, true
			}
		}
	}

	return nil, false
}

func TestGammTestSuite(t *testing.T) {
	suite.Run(t, new(GammTestSuite))
}
