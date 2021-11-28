package types

func NewGenesisState(epoch Epoch) *GenesisState {
	return &GenesisState{Epoch: epoch}
}

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	epoch := Epoch{
		EndBlock: 0,
		Number:   0,
	}
	// epochs := []Epoch{
	// 	{
	// 		Identifier:              "week",
	// 		StartTime:               time.Time{},
	// 		Duration:                time.Hour * 24 * 7,
	// 		CurrentEpoch:            0,
	// 		CurrentEpochStartHeight: 0,
	// 		CurrentEpochStartTime:   time.Time{},
	// 		EpochCountingStarted:    false,
	// 	},
	// }
	return NewGenesisState(epoch)
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	return nil
}
