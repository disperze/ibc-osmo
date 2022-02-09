# InterGAMM

This is a middleware module around [ICS20 transfer module](https://github.com/cosmos/ibc-go/tree/v2.0.3/modules/apps/transfer) to allow executing custom transactions (ex: swap) on osmosis chain.

## IBC Swap

### IbcPacketData

Data packet sent by a blockchain to Osmosis chain to make custom transaction. It contains the following parameters:

| Parameter      | Type      |              Description                          |
| -------------- | --------- | ------------------------------------------------- |
| Denom          | string    | From ICS20                                        | 
| Amount         | string    | From ICS20                                        | 
| Sender         | string    | From ICS20                                        |
| Receiver       | string    | From ICS20                                        |
| gamm           | object    | Osmosis action (optional)                         |

GAMM actions supported: 

### SwapExactAmountInPacketData

Allows to make swap transaction on Osmosis chain, uses `TokenIn` from ICS20 data.

| Parameter         | Type          |              Description                          |
| ----------------- | ------------- | ------------------------------------------------- |
| Sender            | string        | Sender, useful for the caller                     |
| Routes            | [SwapAmountInRoute](https://github.com/osmosis-labs/osmosis/blob/v6.2.0/proto/osmosis/gamm/v1beta1/tx.proto#L81)  | From osmosis                                      |
| TokenOutMinAmount | string        | Min output amount                                 |


## SwapExactAmountInAck

Returns the swap result

| Parameter      | Type      | 
| -------------- | --------- | 
| Denom          | string    |
| Amount         | string    |