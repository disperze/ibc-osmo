# InterGAMM

IBC Implementation for osmosis `x/gamm` module.

## IBC SpotPrice

### SpotPricePacketData

Data packet sent by a blockchain to Osmosis chain to query the spot price of a pool asset. It contains the following parameters:

| Parameter      | Type      |
| -------------- | --------- |
| PoolID         | unit64    |
| TokenIn        | string    |                                              
| TokenOut       | string    |                         

### SpotPricePacketAck

Returns the spot price result

| Parameter      | Type      | 
| -------------- | --------- | 
| Price          | string    |
