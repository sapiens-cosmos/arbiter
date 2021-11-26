package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sapiens-cosmos/arbiter/x/bond/types"
)

// GetBondState returns the bond state.
// If the matched bond policy was not registered to the param store, return the error.
// If the matched bond policy is existent but there is no bond state,
// that case means that it is the first time after the bond policy registered to the param store.
// So, in that case, this returns the initialized bond state.
func (k Keeper) GetBondState(ctx sdk.Context, bondDenom string) (types.BondState, error) {
	_, err := k.GetBondPolicy(ctx, bondDenom)
	if err != nil {
		return types.BondState{}, err
	}

	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetBondStateKey(bondDenom))
	if len(bz) == 0 {
		// If the matched bond policy is existent but there is no bond state,
		// that case means that it is the first time after the bond policy registered to the param store.
		// So, in that case, this returns the initialized bond state.
		return types.BondState{
			TotalDebt:       sdk.NewDec(0),
			LastDecayHeight: 0,
		}, nil
	}

	var bondState types.BondState
	err = k.cdc.UnmarshalBinaryBare(bz, &bondState)
	return bondState, err
}

func (k Keeper) setBondState(ctx sdk.Context, bondDenom string, bondState types.BondState) error {
	store := ctx.KVStore(k.storeKey)

	bz, err := k.cdc.MarshalBinaryBare(&bondState)
	if err != nil {
		return err
	}

	store.Set(types.GetBondStateKey(bondDenom), bz)
	return nil
}
