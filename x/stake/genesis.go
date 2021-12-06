package stake

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sapiens-cosmos/arbiter/x/stake/keeper"
	"github.com/sapiens-cosmos/arbiter/x/stake/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	initTokens := sdk.NewCoins(genState.ModuleAccountBalance, genState.ModuleAccountSTokenBalance)
	err := k.CreateModuleAccount(ctx, initTokens)
	if err != nil {
		panic(err)
	}
	k.SetStakeState(ctx, genState.StakeState)
	k.SetParams(ctx, genState.Params)
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.StakeState = k.GetStakeState(ctx)
	genesis.ModuleAccountBalance = k.GetModuleAccountBalance(ctx)
	genesis.ModuleAccountSTokenBalance = k.GetModuleAccountSTokenBalance(ctx)
	return genesis
}
