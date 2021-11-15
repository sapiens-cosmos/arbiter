package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sapiens-cosmos/arbiter/x/bond/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	// TODO: Set params.

	k.setBondDenom(ctx, genState.BondDenom)
}
