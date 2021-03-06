package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	appParams "github.com/sapiens-cosmos/arbiter/app/params"
	"github.com/sapiens-cosmos/arbiter/x/stake/types"
)

// GetModuleAccountBalance gets the stake coin balance of module account
func (k Keeper) GetModuleAccountAddress(ctx sdk.Context) sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.ModuleName)
}

// GetModuleAccountBalance gets the baseToken balance of module account
func (k Keeper) GetModuleAccountBalance(ctx sdk.Context) sdk.Coin {
	moduleAccAddr := k.GetModuleAccountAddress(ctx)
	return k.bankKeeper.GetBalance(ctx, moduleAccAddr, appParams.BaseCoinUnit)
}

// GetModuleAccountSTokenBalance gets the baseToken balance of module account
func (k Keeper) GetModuleAccountSTokenBalance(ctx sdk.Context) sdk.Coin {
	moduleAccAddr := k.GetModuleAccountAddress(ctx)
	return k.bankKeeper.GetBalance(ctx, moduleAccAddr, appParams.BaseStakeCoinUnit)
}

func (k Keeper) GetAccountSTokenBalance(ctx sdk.Context, addr sdk.AccAddress) sdk.Coin {
	return k.bankKeeper.GetBalance(ctx, addr, appParams.BaseStakeCoinUnit)
}

// CreateModuleAccount creates module account with baseToken and sToken minted
func (k Keeper) CreateModuleAccount(ctx sdk.Context, coins sdk.Coins) error {
	moduleAcc := authtypes.NewEmptyModuleAccount(types.ModuleName, authtypes.Minter, authtypes.Burner)
	k.accountKeeper.SetModuleAccount(ctx, moduleAcc)
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return err
	}
	return nil
}

// RebaseToken increases sToken supply to increase staking balances relative to profit
// individual epoch distribution is equal to profit
func (k Keeper) RebaseToken(ctx sdk.Context, profit int64, epoch int64) sdk.Int {
	rebaseAmount := sdk.NewInt(0)
	totalSupply := k.bankKeeper.GetSupply(ctx).GetTotal().AmountOf(appParams.BaseCoinUnit)

	circulatingSupply := k.CirculatingSupply(ctx)

	// no rebase is done when profit is 0
	if profit == 0 {
		return totalSupply
	} else if circulatingSupply.GT(sdk.ZeroInt()) {
		rebaseAmount = sdk.NewInt(profit).Mul(totalSupply).Quo(circulatingSupply)
	} else {
		rebaseAmount = sdk.NewInt(profit)
	}

	totalSupply.Add(rebaseAmount)
	return totalSupply
}

// CirculatingSupply calculates the circulating supply which is represented through,
// total supply - staked amount
func (k Keeper) CirculatingSupply(ctx sdk.Context) sdk.Int {
	totalSupply := k.bankKeeper.GetSupply(ctx).GetTotal().AmountOf(appParams.BaseCoinUnit)
	moduleAccountBalance := k.GetModuleAccountBalance(ctx).Amount
	return totalSupply.Sub(moduleAccountBalance)
}
