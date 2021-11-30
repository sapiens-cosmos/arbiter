package stake

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sapiens-cosmos/arbiter/x/stake/keeper"
	"github.com/sapiens-cosmos/arbiter/x/stake/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	initTokens := sdk.NewCoins(genState.ModuleAccountBalance, genState.ModuleAccountSTokenBalance)
	k.CreateModuleAccount(ctx, initTokens)
	k.SetEpoch(ctx, genState.Epoch)
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Epoch = k.GetEpoch(ctx)
	genesis.ModuleAccountBalance = k.GetModuleAccountBalance(ctx)
	genesis.ModuleAccountSTokenBalance = k.GetModuleAccountSTokenBalance(ctx)
	return genesis
}
