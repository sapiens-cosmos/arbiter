package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	appParams "github.com/sapiens-cosmos/arbiter/app/params"
	"github.com/sapiens-cosmos/arbiter/x/stake/types"
)

func (k Keeper) JoinStake(ctx sdk.Context, address string, tokenIn sdk.Coin) {
	k.Rebase(ctx)
	k.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.AccAddress(address), types.ModuleName, sdk.Coins{tokenIn})
	k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.Coins{sdk.NewCoin(appParams.BaseStakeCoinUnit, tokenIn.Amount)})

	lock := types.Lock{
		Owner: address,
		Coin:  sdk.NewCoin(appParams.BaseStakeCoinUnit, tokenIn.Amount),
	}
	k.SetLockByAddress(ctx, lock)
}

func (k Keeper) Rebase(ctx sdk.Context) {
	epoch := k.GetEpoch(ctx)
	if epoch.EndBlock < ctx.BlockHeight() {
		k.RebaseToken(ctx, epoch.Distribute, epoch.Number)

		epoch.EndBlock += epoch.Length
		epoch.Number++

		k.Distribute(ctx)

		moduleAccountBalance := k.GetModuleAccountBalance(ctx)
		staked := k.CirculatingSupply(ctx)

		if moduleAccountBalance.Amount.LTE(staked.Amount) {
			epoch.Distribute = 0
		} else {
			epoch.Distribute = moduleAccountBalance.Sub(staked).Amount.Int64()
		}

		k.SetEpoch(ctx, epoch)
	}
}

// Distribute sends epoch rewards to staking contract
func (k Keeper) Distribute(ctx sdk.Context) error {
	totalReward := k.GetTotalReward(ctx)
	moduleAccountSTokenBalance := k.GetModuleAccountSTokenBalance(ctx)

	if totalReward.LT(sdk.NewDecFromInt(k.excessReserves(ctx))) {
		return fmt.Errorf("Insufficient Reserves")
	}

	locks := k.GetAllLocks(ctx)
	// iterate over all locks, calculate the share of an inidividual locks
	// lock only contains perceptual tokens, not the actual sTokens
	// the reward(sToken) relative to share gets minted to module account
	for _, lock := range locks {
		share := lock.Coin.Amount.ToDec().Quo(moduleAccountSTokenBalance.Amount.ToDec())
		reward := totalReward.Mul(share).TruncateInt()
		rewardInBaseToken := sdk.NewCoin(appParams.BaseCoinUnit, reward)
		rewardInSToken := sdk.NewCoin(appParams.BaseStakeCoinUnit, reward)

		err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.Coins{rewardInBaseToken, rewardInSToken})
		if err != nil {
			return err
		}

		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(lock.Owner), sdk.Coins{rewardInSToken})
		if err != nil {
			return err
		}

		lock.Coin = lock.Coin.Add(rewardInSToken)
		k.SetLockByAddress(ctx, lock)
	}
	return nil
}

// Claim receives sTokens from user, unstake and distribute base Token
func (k Keeper) Claim(ctx sdk.Context, address string, amount sdk.Int) error {
	returnCoin := sdk.NewCoin(appParams.BaseStakeCoinUnit, amount)

	// send sTokens from user to module account, then burn
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.AccAddress(address), types.ModuleName, sdk.Coins{returnCoin})
	if err != nil {
		return err
	}
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.Coins{returnCoin})
	if err != nil {
		return err
	}

	// send base Token to user account
	receiveCoin := sdk.NewCoin(appParams.BaseCoinUnit, amount)
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(address), sdk.Coins{receiveCoin})

	// update lock and share status
	lock, err := k.GetLockByAddress(ctx, address)
	if err != nil {
		return err
	}

	remainingStake := lock.Coin.Amount.Sub(amount)

	// if user claimed all stake, delete lock and if not, update lock
	if remainingStake.LTE(sdk.ZeroInt()) {
		k.DeleteLock(ctx, lock)
	} else {
		lock.Coin = lock.Coin.Sub(returnCoin)
		k.SetLockByAddress(ctx, lock)
	}
	return nil
}

func (k Keeper) GetTotalReward(ctx sdk.Context) sdk.Dec {
	totalSupply := k.bankKeeper.GetSupply(ctx, appParams.BaseCoinUnit)
	if totalSupply.Amount.LT(sdk.NewInt(1_000_000)) {
		return sdk.NewDecFromInt(totalSupply.Amount).Quo(sdk.MustNewDecFromStr("0.3058"))
	} else if totalSupply.Amount.GT(sdk.NewInt(1_000_000)) && totalSupply.Amount.LT(sdk.NewInt(10_000_000)) {
		return sdk.NewDecFromInt(totalSupply.Amount).Quo(sdk.MustNewDecFromStr("0.1587"))
	} else if totalSupply.Amount.GT(sdk.NewInt(10_000_000)) && totalSupply.Amount.LT(sdk.NewInt(100_000_000)) {
		return sdk.NewDecFromInt(totalSupply.Amount).Quo(sdk.MustNewDecFromStr("0.1186"))
	} else if totalSupply.Amount.GT(sdk.NewInt(100_000_000)) && totalSupply.Amount.LT(sdk.NewInt(1_000_000_000)) {
		return sdk.NewDecFromInt(totalSupply.Amount).Quo(sdk.MustNewDecFromStr("0.0408"))
	} else {
		return sdk.NewDecFromInt(totalSupply.Amount).Quo(sdk.MustNewDecFromStr("0.0019"))
	}
}

func (k Keeper) AddTotalReserve(ctx sdk.Context, reserve sdk.Int) {
	stakeState := k.GetStakeState(ctx)
	stakeState.TotalReserve = stakeState.TotalReserve.Add(reserve)

	k.SetStakeState(ctx, stakeState)
}

func (k Keeper) GetStakeState(ctx sdk.Context) types.StakeState {
	stakeState := types.StakeState{}
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyStakeState)
	if bz == nil {
		return types.StakeState{}
	}
	err := proto.Unmarshal(bz, &stakeState)
	if err != nil {
		panic(err)
	}
	return stakeState
}

func (k Keeper) SetStakeState(ctx sdk.Context, stakeState types.StakeState) {
	store := ctx.KVStore(k.storeKey)
	value, err := proto.Marshal(&stakeState)
	if err != nil {
		panic(err)
	}
	store.Set(types.KeyStakeState, value)
}

func (k Keeper) excessReserves(ctx sdk.Context) sdk.Int {
	totalSupply := k.bankKeeper.GetSupply(ctx, appParams.BaseCoinUnit)
	totalReserve := k.GetStakeState(ctx).TotalReserve
	return totalReserve.Sub(totalSupply.Amount)
}
