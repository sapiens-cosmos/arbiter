package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/sapiens-cosmos/arbiter/x/stake/types"
)

func permContains(perms []string, perm string) bool {
	for _, v := range perms {
		if v == perm {
			return true
		}
	}

	return false
}

// keeper of the stake store
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      codec.BinaryMarshaler

	paramSpace paramtypes.Subspace

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	paramstore    paramtypes.Subspace
}

// NewKeeper creates a new stake Keeper instance
func NewKeeper(cdc codec.BinaryMarshaler, paramSpace paramtypes.Subspace, key sdk.StoreKey, ak types.AccountKeeper, bk types.BankKeeper, ps paramtypes.Subspace) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}
	_, perms := ak.GetModuleAddressAndPermissions(types.ModuleName)

	if !permContains(perms, authtypes.Burner) {
		panic(fmt.Sprintf("%s module account should have the burner permission", types.ModuleName))
	}

	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		paramSpace:    paramSpace,
		accountKeeper: ak,
		bankKeeper:    bk,
		paramstore:    ps,
	}
}
