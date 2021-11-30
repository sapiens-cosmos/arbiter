package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	appParams "github.com/sapiens-cosmos/arbiter/app/params"
)

// DefaultGenesis returns the default Capability genesis state
// default epoch length is set to 62,000 blocks, which is around 5 human days
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Epoch: Epoch{
			EndBlock:   0,
			Number:     0,
			Length:     62_000,
			Distribute: 0,
		},
		ModuleAccountBalance: sdk.NewCoin(appParams.BaseCoinUnit, sdk.ZeroInt()),
		// sToken Balance starts with 100 shares to avoid zero exception
		ModuleAccountSTokenBalance: sdk.NewCoin(appParams.BaseStakeCoinUnit, sdk.NewInt(100)),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	return nil
}
