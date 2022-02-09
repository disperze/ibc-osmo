package types

import (
	"fmt"
	"strconv"
	"strings"
)

func GetPoolShareDenom(poolId uint64) string {
	return fmt.Sprintf("gamm/pool/%d", poolId)
}

func GetPoolIdFromShareDenom(denom string) (uint64, error) {
	numberStr := strings.TrimLeft(denom, "gamm/pool/")
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return 0, err
	}

	return uint64(number), nil
}
