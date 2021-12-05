package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/sapiens-cosmos/arbiter/x/stake/types"
)

func (k Keeper) GetEpoch(ctx sdk.Context) types.Epoch {
	epoch := types.Epoch{}
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyEpoch)
	if bz == nil {
		return epoch
	}
	err := proto.Unmarshal(bz, &epoch)
	if err != nil {
		panic(err)
	}
	return epoch
}

func (k Keeper) SetEpoch(ctx sdk.Context, epoch types.Epoch) {
	store := ctx.KVStore(k.storeKey)
	value, err := proto.Marshal(&epoch)
	if err != nil {
		panic(err)
	}
	store.Set(types.KeyEpoch, value)
}
