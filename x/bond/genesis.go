package bond

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.Ge)
