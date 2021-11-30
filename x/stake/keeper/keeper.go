package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/sapiens-cosmos/arbiter/x/stake/types"
)

// keeper of the stake store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.BinaryMarshaler
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	paramstore    paramtypes.Subspace
}

// NewKeeper creates a new stake Keeper instance
func NewKeeper(cdc codec.BinaryMarshaler, key sdk.StoreKey, bk types.BankKeeper, ps paramtypes.Subspace) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		bankKeeper: bk,
		paramstore: ps,
	}
}
