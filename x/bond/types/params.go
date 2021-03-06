package types

import (
	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	KeyPolicies = []byte("Policies")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamTable for staking module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(policies []BondPolicy) Params {
	return Params{
		Policies: policies,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		// TODO: Add validate fn.
		paramtypes.NewParamSetPair(KeyPolicies, &p.Policies, func(value interface{}) error {
			return nil
		}),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		[]BondPolicy{
			{
				BondType:        BondType_RESERVE,
				BondDenom:       "bond",
				ControlVariable: sdk.NewDec(1),
				VestingHeight:   10,
			},
		},
	)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
