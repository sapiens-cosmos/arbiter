package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sapiens-cosmos/arbiter/x/bond/types"
)

func (k Keeper) RedeeambleDebt(ctx sdk.Context, bonder sdk.AccAddress) (sdk.Coin, error) {
	debt, err := k.GetDebt(ctx, bonder)
	if err != nil {
		return sdk.Coin{}, err
	}

	heightSince := ctx.BlockHeight() - debt.LastHeight
	vestedRatio := sdk.NewDec(heightSince).QuoInt64(debt.RemainingHeight)

	if vestedRatio.GTE(sdk.NewDec(1)) {
		return sdk.NewCoin(k.GetBaseDenom(ctx), debt.Amount), nil
	}
	payoutAmount := debt.Amount.ToDec().Mul(vestedRatio).TruncateInt()
	return sdk.NewCoin(k.GetBaseDenom(ctx), payoutAmount), nil
}

func (k Keeper) RedeemDebt(ctx sdk.Context, bonder sdk.AccAddress) error {
	debt, err := k.GetDebt(ctx, bonder)
	if err != nil {
		return err
	}

	heightSince := ctx.BlockHeight() - debt.LastHeight
	vestedRatio := sdk.NewDec(heightSince).QuoInt64(debt.RemainingHeight)

	if vestedRatio.GTE(sdk.NewDec(1)) {
		k.deleteDebt(ctx, bonder)

		return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, bonder, sdk.NewCoins(
			sdk.NewCoin(k.GetBaseDenom(ctx), debt.Amount),
		))
	}
	payoutAmount := debt.Amount.ToDec().Mul(vestedRatio).TruncateInt()

	newDebt := types.Debt{
		Amount:          debt.Amount.Sub(payoutAmount),
		RemainingHeight: debt.RemainingHeight - (ctx.BlockHeight() - debt.LastHeight),
		LastHeight:      ctx.BlockHeight(),
	}

	err = k.setDebt(ctx, newDebt, bonder)
	if err != nil {
		return err
	}

	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, bonder, sdk.NewCoins(
		sdk.NewCoin(k.GetBaseDenom(ctx), payoutAmount),
	))
}

func (k Keeper) AddDebt(ctx sdk.Context, bonder sdk.AccAddress, bondDenom string, amount sdk.Int) error {
	policy, err := k.GetBondPolicy(ctx, bondDenom)
	if err != nil {
		return err
	}

	debt, err := k.GetDebt(ctx, bonder)
	if err != nil {
		if errors.Is(err, types.ErrNoDebt) {
			debt = types.Debt{
				Amount: sdk.NewInt(0),
			}
		} else {
			return err
		}
	}

	debt = types.Debt{
		Amount:          debt.Amount.Add(amount),
		RemainingHeight: policy.VestingHeight,
		LastHeight:      ctx.BlockHeight(),
	}

	return k.setDebt(ctx, debt, bonder)
}

func (k Keeper) GetDebt(ctx sdk.Context, bonder sdk.AccAddress) (types.Debt, error) {
	store := ctx.KVStore(k.storeKey)

	key := types.GetDebtKey(bonder)

	bz := store.Get(key)
	if len(bz) == 0 {
		return types.Debt{}, sdkerrors.Wrap(types.ErrNoDebt, bonder.String())
	}

	var debt types.Debt
	err := k.cdc.UnmarshalBinaryBare(bz, &debt)
	return debt, err
}

func (k Keeper) setDebt(ctx sdk.Context, debt types.Debt, bonder sdk.AccAddress) error {
	store := ctx.KVStore(k.storeKey)

	key := types.GetDebtKey(bonder)

	bz, err := k.cdc.MarshalBinaryBare(&debt)
	if err != nil {
		return err
	}

	store.Set(key, bz)
	return nil
}

func (k Keeper) deleteDebt(ctx sdk.Context, bonder sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	key := types.GetDebtKey(bonder)

	store.Delete(key)
}

// GetDebtRatio returns the bond's debt ratio.
// Debt Ratio = Bonds Outstanding / base Supply
func (k Keeper) GetDebtRatio(ctx sdk.Context, bondDenom string) (sdk.Dec, error) {
	state, err := k.GetBondState(ctx, bondDenom)
	if err != nil {
		return sdk.Dec{}, err
	}

	supply := k.bankKeeper.GetSupply(ctx).GetTotal().AmountOf(bondDenom)
	return state.TotalDebt.Quo(supply.ToDec()), nil
}

// inflateTotalDebt inflates the amount of debt to the bond state.
func (k Keeper) inflateTotalDebt(ctx sdk.Context, bondDenom string, amount sdk.Int) error {
	state, err := k.GetBondState(ctx, bondDenom)
	if err != nil {
		return err
	}

	state.TotalDebt = state.TotalDebt.Add(amount.ToDec())
	return k.setBondState(ctx, bondDenom, state)
}

// DecayDebt decays debt linearly.
func (k Keeper) DecayDebt(ctx sdk.Context, bondDenom string) error {
	state, err := k.GetBondState(ctx, bondDenom)
	if err != nil {
		return err
	}

	policy, err := k.GetBondPolicy(ctx, bondDenom)
	if err != nil {
		return err
	}

	heightSince := ctx.BlockHeight() - state.LastDecayHeight

	if heightSince == 0 {
		return nil
	}

	if heightSince >= policy.VestingHeight {
		state.TotalDebt = sdk.NewDec(0)
	} else {
		state.TotalDebt = state.TotalDebt.MulInt64(policy.VestingHeight - heightSince).QuoInt64(policy.VestingHeight)
	}

	state.LastDecayHeight = ctx.BlockHeight()

	return k.setBondState(ctx, bondDenom, state)
}

// DecayAllDebt decays all debt.
func (k Keeper) DecayAllDebt(ctx sdk.Context) error {
	policies := k.GetBondPolicies(ctx)

	for _, policy := range policies {
		err := k.DecayDebt(ctx, policy.BondDenom)
		if err != nil {
			return err
		}
	}

	return nil
}
