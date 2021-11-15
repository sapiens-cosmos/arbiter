package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/sapiens-cosmos/arbiter/x/bond/types"
)

func (k Keeper) deposit(amount sdk.Int, maxPrice sdk.Dec, depositor string) error {
	address, err := sdk.AccAddressFromBech32(depositor)
	if err != nil {
		return err
	}
	k.decayDebt(ctx)
}

//
func (k Keeper) decayDebt(ctx sdk.Context) {
	bondState := k.GetBondState(ctx)
	terms := k.GetTerms(ctx)
	totalDebt := bondState.TotalDebt

	blockSinceDecay := ctx.BlockHeight() - bondState.LastDecay
	decay := totalDebt.MulInt64(blockSinceDecay).QuoInt64(terms.VestingTerm)
	if decay.GT(totalDebt) {
		decay = totalDebt
	}

	totalDebt = totalDebt.Sub(decay)

	bondState.TotalDebt = totalDebt
	bondState.LastDecay = ctx.BlockHeight()

	k.SetBondState(ctx, bondState)
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
