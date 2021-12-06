package keeper

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	appParams "github.com/sapiens-cosmos/arbiter/app/params"
	"github.com/sapiens-cosmos/arbiter/x/stake/types"
)

func (k Keeper) JoinStake(ctx sdk.Context, address string, tokenIn sdk.Coin) error {
	err := k.Rebase(ctx)
	if err != nil {
		return err
	}
	accAddress, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, accAddress, types.ModuleName, sdk.Coins{tokenIn})
	if err != nil {
		return err
	}
	sTokens := sdk.NewCoin(appParams.BaseStakeCoinUnit, tokenIn.Amount)
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.Coins{sTokens})
	if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, accAddress, sdk.Coins{sTokens})
	if err != nil {
		return err
	}

	lock, err := k.GetLockByAddress(ctx, address)
	if err != nil {
		if errors.Is(err, types.ErrNoLock) {
			lock = types.Lock{
				Owner: address,
				Coin:  sTokens,
			}

			k.SetLockByAddress(ctx, lock)
			return nil
		}
		return err
	}

	lock.Coin = lock.Coin.Add(sTokens)
	k.SetLockByAddress(ctx, lock)
	return nil
}

func (k Keeper) Rebase(ctx sdk.Context) error {
	stakeState := k.GetStakeState(ctx)

	if stakeState.Epoch.EndBlock < ctx.BlockHeight() {
		k.RebaseToken(ctx, stakeState.Epoch.Distribute, stakeState.Epoch.Number)

		stakeState.Epoch.EndBlock += stakeState.Epoch.Length
		stakeState.Epoch.Number++

		err := k.Distribute(ctx)
		if err != nil {
			return err
		}

		moduleAccountBalance := k.GetModuleAccountBalance(ctx)
		staked := k.CirculatingSupply(ctx)

		if moduleAccountBalance.Amount.LTE(staked) {
			stakeState.Epoch.Distribute = 0
		} else {
			stakeState.Epoch.Distribute = moduleAccountBalance.Amount.Sub(staked).Int64()
		}
		k.SetStakeState(ctx, stakeState)
	}
	return nil
}

// Distribute sends epoch rewards to staking contract
func (k Keeper) Distribute(ctx sdk.Context) error {
	totalReward := k.GetTotalReward(ctx)
	moduleAccountSTokenBalance := k.GetModuleAccountSTokenBalance(ctx)

	if totalReward.LT(sdk.NewDecFromInt(k.excessReserves(ctx))) {
		return fmt.Errorf("insufficient Reserves")
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

		accAddress, err := sdk.AccAddressFromBech32(lock.Owner)
		if err != nil {
			return err
		}
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, accAddress, sdk.Coins{rewardInSToken})
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
	accAddress, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, accAddress, types.ModuleName, sdk.Coins{returnCoin})
	if err != nil {
		return err
	}
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.Coins{returnCoin})
	if err != nil {
		return err
	}

	// send base Token to user account
	receiveCoin := sdk.NewCoin(appParams.BaseCoinUnit, amount)
	accAddress, err = sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, accAddress, sdk.Coins{receiveCoin})
	if err != nil {
		return err
	}

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

func (k Keeper) GetBalance(ctx sdk.Context, address sdk.AccAddress) sdk.Int {
	return k.bankKeeper.GetBalance(ctx, address, appParams.BaseCoinUnit).Amount
}

func (k Keeper) GetRewardYield(ctx sdk.Context) sdk.Dec {
	totalStaked := k.GetModuleAccountBalance(ctx)
	distributed := k.GetTotalReward(ctx)
	if totalStaked.Amount.Equal(sdk.ZeroInt()) || distributed.Equal(sdk.ZeroDec()) {
		return sdk.ZeroDec()
	}
	rewardYield := distributed.Quo(totalStaked.Amount.ToDec())
	return rewardYield
}

func (k Keeper) GetTotalReward(ctx sdk.Context) sdk.Dec {
	totalSupply := k.bankKeeper.GetSupply(ctx).GetTotal().AmountOf(appParams.BaseCoinUnit)
	rewardRate := k.GetParams(ctx).RewardRate
	return sdk.NewDecFromInt(totalSupply).Mul(rewardRate)
}

func (k Keeper) GetBlockUntilRebase(ctx sdk.Context) int64 {
	epoch := k.GetStakeState(ctx).Epoch
	return epoch.EndBlock - ctx.BlockHeight()
}

func (k Keeper) AddTotalReserve(ctx sdk.Context, reserve sdk.Int) {
	stakeState := k.GetStakeState(ctx)
	stakeState.TotalReserve = stakeState.TotalReserve.Add(reserve)

	k.SetStakeState(ctx, stakeState)
}

func (k Keeper) GetTotalReserve(ctx sdk.Context) sdk.Int {
	stakeState := k.GetStakeState(ctx)

	return stakeState.TotalReserve
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
	totalSupply := k.bankKeeper.GetSupply(ctx).GetTotal().AmountOf(appParams.BaseCoinUnit)
	totalReserve := k.GetStakeState(ctx).TotalReserve
	return totalReserve.Sub(totalSupply)
}
