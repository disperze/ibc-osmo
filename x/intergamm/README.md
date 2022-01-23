# InterGAMM

IBC Implementation for osmosis `x/gamm` module.

## IBC SpotPrice

### SpotPricePacketData

Data packet sent by a blockchain to Osmosis chain to query the spot price of a pool asset. It contains the following parameters:

| Parameter      | Type      |              Description                   |
| -------------- | --------- | ------------------------------------------ |
| ClientID       | string    | Unique identifier specified by the client. |
| PoolID         | unit64    | Pool asset ID                              |
| TokenIn        | string    | Token In denom                             |                          
| TokenOut       | string    | Token Out denom                            |               

### SpotPricePacketAck

Returns the spot price result

| Parameter      | Type      | 
| -------------- | --------- | 
| Price          | string    |
