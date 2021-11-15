package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/sapiens-cosmos/arbiter/x/bond/types"
)

func (k Keeper) deposit(ctx sdk.Context, amount sdk.Int, maxPrice sdk.Dec, depositor string) error {
	bondState := k.GetBondState(ctx)
	address, err := sdk.AccAddressFromBech32(depositor)
	if err != nil {
		return err
	}
	k.DecayDebt(ctx)

	priceInUSD := k.BondPriceInUSD(ctx)
	nativePrice := k.BondPrice(ctx)

	value := k.Valuation(bondState.Principle, amount)
	payOut := value.Quo(k.BondPrice(ctx))
}

//
func (k Keeper) DecayDebt(ctx sdk.Context) {
	bondState := k.GetBondState(ctx)
	totalDebt := bondState.TotalDebt

	decay := k.DebtToDecay(ctx)

	totalDebt = totalDebt.Sub(decay)

	bondState.TotalDebt = totalDebt
	bondState.LastDecay = ctx.BlockHeight()

	k.SetBondState(ctx, bondState)
}

func (k Keeper) DebtToDecay(ctx sdk.Context) sdk.Dec {
	bondState := k.GetBondState(ctx)
	terms := k.GetTerms(ctx)

	totalDebt := bondState.TotalDebt
	blockSinceDecay := ctx.BlockHeight() - bondState.LastDecay
	decay := totalDebt.MulInt64(blockSinceDecay).QuoInt64(terms.VestingTerm)
	if decay.GT(totalDebt) {
		decay = totalDebt
	}
	return decay
}

// calculate surrent ratio of debt to bondDenom supply
func (k Keeper) DebtRatio(ctx sdk.Context) sdk.Dec {
	bondDenom := k.GetBondDenom(ctx)

	totalSupply := k.bankKeeper.GetSupply(ctx).GetTotal().AmountOf(bondDenom)

	debtRatio := k.CurrentDebt(ctx).QuoInt(totalSupply)
	return debtRatio
}

func (k Keeper) BondPrice(ctx sdk.Context) sdk.Dec {
	terms := k.GetTerms(ctx)

	price := k.DebtRatio(ctx).Mul(terms.ControlVariable)
	return price
}

// BondPriceInUSD converts bond price to USD
func (k Keeper) BondPriceInUSD(ctx sdk.Context) sdk.Dec {
	bondState := k.GetBondState(ctx)

	if bondState.IsLiquidityBond {
		return k.DebtRatio(ctx).Mul(k.LiquidityPairToBondDenom()).QuoInt64(100)
	} else {
		return k.DebtRatio(ctx).QuoInt64(100)
	}
}

func (k Keeper) CurrentDebt(ctx sdk.Context) sdk.Dec {
	totalDebt := k.GetBondState(ctx).TotalDebt
	debtToDecay := k.DebtToDecay(ctx)
	return totalDebt.Sub(debtToDecay)
}

func (k Keeper) GetBondState(ctx sdk.Context) types.BondState {
	bondState := types.BondState{}
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.KeyBondState)
	if b == nil {
		return bondState
	}
	err := proto.Unmarshal(b, &bondState)
	if err != nil {
		return bondState
	}
}

func (k Keeper) SetBondState(ctx sdk.Context, bondState types.BondState) {
	store := ctx.KVStore(k.storeKey)
	bondStateKey := types.KeyBondState

	value, err := proto.Marshal(&bondState)
	if err != nil {
		panic(err)
	}

	store.Set(bondStateKey, value)
}

func (k Keeper) GetTerms(ctx sdk.Context) types.Terms {
	terms := types.Terms{}
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.KeyTerms)
	if b == nil {
		return terms
	}
	err := proto.Unmarshal(b, &terms)
	if err != nil {
		return terms
	}
}

func (k Keeper) SetTerms(ctx sdk.Context, terms types.Terms) {
	store := ctx.KVStore(k.storeKey)
	termsKey := types.KeyTerms

	value, err := proto.Marshal(&terms)
	if err != nil {
		panic(err)
	}

	store.Set(termsKey, value)
}

func (k Keeper) GetBondDenom(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyBondDenom)
	if len(bz) == 0 {
		panic("bond denom not set")
	}

	return string(bz)
}

func (k Keeper) setBondDenom(ctx sdk.Context, bondDenom string) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.KeyBondDenom, []byte(bondDenom))
}
