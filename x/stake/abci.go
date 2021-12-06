package stake

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sapiens-cosmos/arbiter/x/stake/keeper"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	k.Rebase(ctx)
}
