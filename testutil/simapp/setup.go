package simapp

import (
	"encoding/json"

	ibctesting "github.com/cosmos/ibc-go/v2/testing"
	ibcsimapp "github.com/cosmos/ibc-go/v2/testing/simapp"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

func SetupTestingApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
	db := dbm.NewMemDB()
	encCdc := MakeTestEncodingConfig()
	app := NewSimApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 5, encCdc, ibcsimapp.EmptyAppOptions{})
	return app, NewDefaultGenesisState(encCdc.Marshaler)
}
