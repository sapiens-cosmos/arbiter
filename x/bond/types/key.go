package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName is the name of the bond module
	ModuleName = "bond"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// RouterKey is the msg router key for the staking module
	RouterKey = ModuleName
)

var (
	// KeyBaseDenom defines state to store base denom
	KeyBaseDenom = []byte{0x01}
	// KeyBondState defines state to store bond state
	KeyBondState = []byte{0x02}
	// KeyDebt defines state to store debt
	KeyDebt = []byte{0x03}
)

func GetBondStateKey(bondDenom string) []byte {
	return append(KeyBondState, []byte(bondDenom)...)
}

func GetDebtKey(bonder sdk.AccAddress) []byte {
	return append(KeyDebt, bonder...)
}
