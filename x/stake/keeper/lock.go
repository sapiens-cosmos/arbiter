package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"

	"github.com/sapiens-cosmos/arbiter/x/stake/types"
)

func (k Keeper) GetAllLocks(ctx sdk.Context) []types.Lock {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, []byte(types.KeyPrefixLock))
	iterator := prefixStore.Iterator(nil, nil)
	locks := []types.Lock{}

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		lock := types.Lock{}
		err := proto.Unmarshal(iterator.Value(), &lock)
		if err != nil {
			panic(err)
		}
		locks = append(locks, lock)
	}
	return locks
}

func (k Keeper) GetLockByAddress(ctx sdk.Context, address string) (types.Lock, error) {
	store := ctx.KVStore(k.storeKey)
	lockKey := types.GetAddressLockStoreKey(address)
	if !store.Has(lockKey) {
		return types.Lock{}, fmt.Errorf("lock of address %s does not exist", address)
	}

	bz := store.Get(lockKey)
	lock := types.Lock{}
	err := proto.Unmarshal(bz, &lock)
	if err != nil {
		return lock, err
	}
	return lock, nil
}

func (k Keeper) SetLockByAddress(ctx sdk.Context, lock types.Lock) {
	store := ctx.KVStore(k.storeKey)
	value, err := proto.Marshal(&lock)
	if err != nil {
		panic(err)
	}
	lockKey := types.GetAddressLockStoreKey(lock.Owner)
	store.Set(lockKey, value)
}

func (k Keeper) DeleteLock(ctx sdk.Context, lock types.Lock) {
	store := ctx.KVStore(k.storeKey)
	lockKey := types.GetAddressLockStoreKey(lock.Owner)
	store.Delete(lockKey)
}

// func (k Keeper) GetLockByID(ctx sdk.Context, lockID uint64) (*types.Lock, error) {
// 	lock := types.Lock{}
// 	store := ctx.KVStore(k.storeKey)
// 	lockKey := lockStoreKey(lockID)
// 	if !store.Has(lockKey) {
// 		return nil, fmt.Errorf("lock with ID %d does not exist", lockID)
// 	}
// 	bz := store.Get(lockKey)
// 	err := proto.Unmarshal(bz, &lock)
// 	return &lock, err
// }

// // GetLastLockID returns ID used last time
// func (k Keeper) GetLastLockID(ctx sdk.Context) uint64 {
// 	store := ctx.KVStore(k.storeKey)
// 	bz := store.Get(types.KeyLastLockID)
// 	if bz == nil {
// 		return 0
// 	}

// 	return sdk.BigEndianToUint64(bz)
// }

// // SetLastLockID save ID used by last lock
// func (k Keeper) SetLastLockID(ctx sdk.Context, ID uint64) {
// 	store := ctx.KVStore(k.storeKey)
// 	store.Set(types.KeyLastLockID, sdk.Uint64ToBigEndian(ID))
// }

// // lockStoreKey returns action store key from ID
// func lockStoreKey(ID uint64) []byte {
// 	return combineKeys(types.KeyPrefixLock, sdk.Uint64ToBigEndian(ID))
// }

// // combineKeys combine bytes array into a single bytes
// func combineKeys(keys ...[]byte) []byte {
// 	return bytes.Join(keys, types.KeyIndexSeparator)
// }
