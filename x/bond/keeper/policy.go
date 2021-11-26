package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sapiens-cosmos/arbiter/x/bond/types"
)

// GetBondPolicies returns the bond policies from param store.
func (k Keeper) GetBondPolicies(ctx sdk.Context) (res []types.BondPolicy) {
	k.paramstore.Get(ctx, types.KeyPolicies, &res)
	return
}

// GetBondPolicy finds the matched bond policy from param store.
// If there is no matched bond policy, return the error.
func (k Keeper) GetBondPolicy(ctx sdk.Context, bondDenom string) (types.BondPolicy, error) {
	policies := k.GetBondPolicies(ctx)

	for _, existingPolicy := range policies {
		if existingPolicy.BondDenom == bondDenom {
			return existingPolicy, nil
		}
	}

	return types.BondPolicy{}, sdkerrors.Wrap(types.ErrNoBondPolicy, bondDenom)
}
