package intergamm_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v2/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v2/modules/core/04-channel/types"

	"github.com/disperze/ibc-osmo/x/intergamm/types"
)

// constructs a send from chainA to chainB on the established channel/connection
// and sends the same coin back from chainB to chainA.
func (suite *GammTestSuite) TestOnRecvJoinPoolPacket() {
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
			"invalid min share amount", func() {
				joinPoolPacket := types.NewIbcJoinPoolPacketData(
					suite.chainA.SenderAccount.GetAddress().String(),
					suite.chainB.SenderAccount.GetAddress().String(),
					"1000000",
					sdk.DefaultBondDenom,
					1,
					sdk.NewInt(-1),
				)
				packetData = joinPoolPacket.GetBytes()
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

			amountToSend, ok := sdk.NewIntFromString("1000000")
			suite.Require().True(ok)
			coinToSendToB := sdk.NewCoin(sdk.DefaultBondDenom, amountToSend)
			expDenomOut := "gamm/pool/1"
			expAmountOut := "2000000"
			amountSharesOut, ok := sdk.NewIntFromString(expAmountOut)
			suite.Require().True(ok)

			tokenPacket := types.NewIbcJoinPoolPacketData(
				suite.chainA.SenderAccount.GetAddress().String(),
				suite.chainB.SenderAccount.GetAddress().String(),
				coinToSendToB.Amount.String(),
				coinToSendToB.Denom,
				1,
				amountSharesOut,
			)
			packetData = tokenPacket.GetBytes()

			tc.malleate()

			// Send packet
			packet := channeltypes.NewPacket(packetData, seq, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
			err := packet.ValidateBasic()
			suite.Require().NoError(err)

			previousSeqB, ok := suite.GetSimApp(suite.chainB).IBCKeeper.ChannelKeeper.GetNextSequenceSend(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
			suite.Require().True(ok)

			// fund address balance is empty
			fundAddress := types.GetFundAddress(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
			balance := suite.GetSimApp(suite.chainB).BankKeeper.GetBalance(suite.chainB.GetContext(), fundAddress, expDenomOut)
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
			balance = suite.GetSimApp(suite.chainB).BankKeeper.GetBalance(suite.chainB.GetContext(), fundAddress, expDenomOut)
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
				var shares types.JoinPoolAck
				err = types.ModuleCdc.UnmarshalJSON(res.Result, &shares)
				suite.Require().NoError(err)

				suite.Require().Equal(expAmountOut, shares.Amount)
				suite.Require().Equal(expDenomOut, shares.Denom)
			}
		})
	}
}

func (suite *GammTestSuite) TestOnRecvExitPoolPacket() {
	var packetData []byte
	testCases := []struct {
		name          string
		malleate      func()
		expAckSuccess bool
	}{
		{
			"success", func() {}, false,
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

			amountToSend, ok := sdk.NewIntFromString("1000000")
			suite.Require().True(ok)
			// denomPrefix := transfertypes.GetDenomPrefix(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
			// TODO: denom with slash prefix over IBC, ex: "gamm/pool/1"
			coinToSendToB := sdk.NewCoin(sdk.DefaultBondDenom, amountToSend)
			expDenomOut := "uosmo"
			expAmountOut := "2000000"
			amountMintOut, ok := sdk.NewIntFromString(expAmountOut)
			suite.Require().True(ok)

			tokenPacket := types.NewIbcExitPoolPacketData(
				suite.chainA.SenderAccount.GetAddress().String(),
				suite.chainB.SenderAccount.GetAddress().String(),
				coinToSendToB.Amount.String(),
				coinToSendToB.Denom,
				expDenomOut,
				amountMintOut,
			)
			packetData = tokenPacket.GetBytes()

			tc.malleate()

			// Send packet
			packet := channeltypes.NewPacket(packetData, seq, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
			err := packet.ValidateBasic()
			suite.Require().NoError(err)

			previousSeqB, ok := suite.GetSimApp(suite.chainB).IBCKeeper.ChannelKeeper.GetNextSequenceSend(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
			suite.Require().True(ok)

			// fund address balance is empty
			fundAddress := types.GetFundAddress(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
			balance := suite.GetSimApp(suite.chainB).BankKeeper.GetBalance(suite.chainB.GetContext(), fundAddress, expDenomOut)
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
			balance = suite.GetSimApp(suite.chainB).BankKeeper.GetBalance(suite.chainB.GetContext(), fundAddress, expDenomOut)
			suite.Require().Equal(sdk.ZeroInt(), balance.Amount)

			ctx := suite.chainB.GetContext()
			_, found := path.EndpointB.Chain.App.GetIBCKeeper().ChannelKeeper.GetPacketAcknowledgement(ctx, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, seq)
			suite.Require().True(found)

			// check packet acknowledgement
			var ack channeltypes.Acknowledgement
			err = types.ModuleCdc.UnmarshalJSON(ackBytes, &ack)
			suite.Require().NoError(err)

			// suite.T().Logf("ack: %+v", string(ackBytes))
			res, isSuccess := ack.Response.(*channeltypes.Acknowledgement_Result)
			suite.Require().Equal(tc.expAckSuccess, isSuccess)

			if tc.expAckSuccess {
				var shares types.ExitPoolAck
				err = types.ModuleCdc.UnmarshalJSON(res.Result, &shares)
				suite.Require().NoError(err)

				suite.Require().Equal(expAmountOut, shares.Amount)
				suite.Require().Equal(expDenomOut, shares.Denom)
			}
		})
	}
}
