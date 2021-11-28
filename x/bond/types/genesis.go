package types

// NewGenesisState creates a new GenesisState instanc
func NewGenesisState(params Params, baseDenom string) *GenesisState {
	return &GenesisState{
		Params:    params,
		BaseDenom: baseDenom,
	}
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), "stake")
}

func (genesis GenesisState) Validate() error {
	// TODO
	return nil
}
