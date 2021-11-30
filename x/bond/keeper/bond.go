package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sapiens-cosmos/arbiter/x/bond/types"
)

// BondIn bonds the amount of coin to receive the base coin.
func (k Keeper) BondIn(ctx sdk.Context, bonder sdk.AccAddress, coin sdk.Coin) error {
	premium, err := k.GetPremium(ctx, coin.Denom)
	if err != nil {
		return err
	}

	riskFreePrice, err := k.GetRiskFreePrice(ctx, coin.Denom)
	if err != nil {
		return err
	}

	// Executing Price = RiskFreePrice * Premium {Premium ≥ 1}
	executingPrice := riskFreePrice.ToDec().Mul(premium).TruncateInt()

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, bonder, types.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		return err
	}

	debtAmount := coin.Amount.Quo(executingPrice)

	err = k.inflateTotalDebt(ctx, coin.Denom, debtAmount)
	if err != nil {
		return err
	}

	// Mint the base coin
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(k.GetBaseDenom(ctx), debtAmount)))
	if err != nil {
		return err
	}

	err = k.AddDebt(ctx, bonder, coin.Denom, debtAmount)
	if err != nil {
		return err
	}

	// profit := debtAmount.Mul(executingPrice.Sub(riskFreePrice))
	panic("add logic to distribute profit to the stakers")
}

// GetPremium returns the premium to calculate the executing price.
// Premium = 1 + (Debt Ratio * ControlVariable)
func (k Keeper) GetPremium(ctx sdk.Context, bondDenom string) (sdk.Dec, error) {
	debtRatio, err := k.GetDebtRatio(ctx, bondDenom)
	if err != nil {
		return sdk.Dec{}, err
	}

	policy, err := k.GetBondPolicy(ctx, bondDenom)
	if err != nil {
		return sdk.Dec{}, err
	}

	return sdk.NewDec(1).Add(debtRatio.Mul(policy.ControlVariable)), nil
}

func (k Keeper) GetRiskFreePrice(ctx sdk.Context, bondDenom string) (sdk.Int, error) {
	policy, err := k.GetBondPolicy(ctx, bondDenom)
	if err != nil {
		return sdk.Int{}, err
	}

	switch policy.BondType {
	case types.BondType_RESERVE:
		// If the bond type is reverse, just use the constant 1.
		return sdk.NewInt(1), nil
	case types.BondType_LIQUIDITY:
		panic("TODO: risk free price for the liquidity not yet implemented")
	default:
		panic("unknown bond type")
	}
}

func (k Keeper) GetBaseDenom(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyBaseDenom)
	if len(bz) == 0 {
		panic("bond denom not set")
	}

	return string(bz)
}

func (k Keeper) setBaseDenom(ctx sdk.Context, bondDenom string) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.KeyBaseDenom, []byte(bondDenom))
}
