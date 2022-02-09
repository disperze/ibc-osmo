package types

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
}

var (
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

// SafeUnmarshalJSON unmarshals JSON data but allow unknown fields.
func SafeUnmarshalJSON(pc *codec.ProtoCodec, bz []byte, ptr proto.Message) error {
	m, ok := ptr.(codec.ProtoMarshaler)
	if !ok {
		return fmt.Errorf("cannot protobuf JSON decode unsupported type: %T", ptr)
	}

	unmarshaler := jsonpb.Unmarshaler{AnyResolver: pc.InterfaceRegistry(), AllowUnknownFields: true}
	err := unmarshaler.Unmarshal(strings.NewReader(string(bz)), m)
	if err != nil {
		return err
	}

	return cdctypes.UnpackInterfaces(ptr, pc.InterfaceRegistry())
}
