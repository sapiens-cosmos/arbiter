package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

// LiquidityPairToBondDenom is markdown for liquidity
func (k Keeper) LiquidityPairToBondDenom() sdk.Dec {
	//reserve1, reserve2 := k.GetReserves()

	panic("TODO: implement oracle")
}

func (k Keeper) GetReserves() (sdk.Dec, sdk.Dec) {
	panic("TODO: implement reserves oracle")
}

// Valuation returns valuation of the lptoken and the amount
func (k Keeper) Valuation(lptoken string, amount sdk.Int) sdk.Dec {
	panic("TODO: implement valuation method")
}
