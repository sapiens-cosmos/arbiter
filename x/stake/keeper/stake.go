package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) JoinStake(ctx sdk.Context) {

}

func (k Keeper) Rebase(ctx sdk.Context) {
	epoch := k.GetEpoch(ctx)
	if epoch.EndBlock < ctx.BlockHeight() {

		k.RebaseToken(ctx, epoch.Distribute, epoch.Number)
		k.Distribute(ctx)

		epoch.EndBlock += epoch.Length
		epoch.Number++

		k.SetEpoch(ctx, epoch)
	}
}

// Distribute sends epoch rewards to staking contract
func (k Keeper) Distribute(ctx sdk.Context) {
	epoch := k.GetEpoch(ctx)
	locks := k.GetAllLocks(ctx)

}

func (k Keeper) GetReward(ctx sdk.Context) {

}
