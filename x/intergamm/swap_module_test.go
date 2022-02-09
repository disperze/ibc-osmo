package intergamm_test

import (
	"fmt"
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
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

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
			"invalid routes", func() {
				swapTokenPacket := types.NewIbcPacketData(
					suite.chainA.SenderAccount.GetAddress().String(),
					suite.chainB.SenderAccount.GetAddress().String(),
					"100",
					sdk.DefaultBondDenom,
					[]types.SwapAmountInRoute{},
					sdk.OneInt(),
				)
				packetData = swapTokenPacket.GetBytes()
			}, false,
		},
		{
			"invalid denom out", func() {
				swapTokenPacket := types.NewIbcPacketData(
					suite.chainA.SenderAccount.GetAddress().String(),
					suite.chainB.SenderAccount.GetAddress().String(),
					"100",
					sdk.DefaultBondDenom,
					[]types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: simapp.InvalidDenom,
						},
					},
					sdk.OneInt(),
				)
				packetData = swapTokenPacket.GetBytes()
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

			amountToSwap, ok := sdk.NewIntFromString("1000000")
			swapDenomOut := "uosmo"
			suite.Require().True(ok)
			coinToSwapToB := sdk.NewCoin(sdk.DefaultBondDenom, amountToSwap)

			swapTokenPacket := types.NewIbcPacketData(
				suite.chainA.SenderAccount.GetAddress().String(),
				suite.chainB.SenderAccount.GetAddress().String(),
				coinToSwapToB.Amount.String(),
				coinToSwapToB.Denom,
				[]types.SwapAmountInRoute{
					{
						PoolId:        1,
						TokenOutDenom: swapDenomOut,
					},
				},
				sdk.OneInt(),
			)
			packetData = swapTokenPacket.GetBytes()

			tc.malleate()

			// Send packet
			packet := channeltypes.NewPacket(packetData, seq, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
			err := packet.ValidateBasic()
			suite.Require().NoError(err)

			previousSeqB, ok := suite.GetSimApp(suite.chainB).IBCKeeper.ChannelKeeper.GetNextSequenceSend(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
			suite.Require().True(ok)

			// fund address balance is empty
			fundAddress := types.GetFundAddress(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
			balance := suite.GetSimApp(suite.chainB).BankKeeper.GetBalance(suite.chainB.GetContext(), fundAddress, swapDenomOut)
			suite.Require().Equal(sdk.ZeroInt(), balance.Amount)

			err = path.EndpointA.SendPacket(packet)
			suite.Require().NoError(err)

			ackBytes, err := suite.RecvPacket(path.EndpointB, packet)
			suite.Require().NoError(err)

			nextSeqB, ok := suite.GetSimApp(suite.chainB).IBCKeeper.ChannelKeeper.GetNextSequenceSend(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
			suite.Require().True(ok)

			// ICS4wrapper avoid to send packet (only used as "safeMint" for output tokens)
			suite.Require().Equal(previousSeqB, nextSeqB)

			// fund address balance is empty after swap
			balance = suite.GetSimApp(suite.chainB).BankKeeper.GetBalance(suite.chainB.GetContext(), fundAddress, swapDenomOut)
			suite.Require().Equal(sdk.ZeroInt(), balance.Amount)

			ctx := suite.chainB.GetContext()
			_, found := path.EndpointB.Chain.App.GetIBCKeeper().ChannelKeeper.GetPacketAcknowledgement(ctx, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, seq)
			suite.Require().True(found)

			// check packet acknowledgement
			var ack channeltypes.Acknowledgement
			err = types.ModuleCdc.UnmarshalJSON(ackBytes, &ack)
			suite.Require().NoError(err)

			res, isSuccess := ack.Response.(*channeltypes.Acknowledgement_Result)
			suite.Require().Equal(tc.expAckSuccess, isSuccess)

			if tc.expAckSuccess {
				var swapOut types.SwapExactAmountInAck
				err = types.ModuleCdc.UnmarshalJSON(res.Result, &swapOut)
				suite.Require().NoError(err)

				suite.Require().NotEmpty(swapOut.Amount)
				suite.Require().NotEmpty(swapOut.Denom)
			}
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

func (suite *GammTestSuite) GetSimApp(chain *ibctesting.TestChain) *simapp.SimApp {
	app, ok := chain.App.(*simapp.SimApp)
	suite.Require().True(ok)

	return app
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
