package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sapiens-cosmos/arbiter/x/bond/types"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServerImpl returns an implementation of the bond MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) BondIn(goCtx context.Context, msg *types.MsgBondIn) (*types.MsgBondInResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	bonder, err := sdk.AccAddressFromBech32(msg.Bonder)
	if err != nil {
		return nil, err
	}

	err = k.keeper.BondIn(ctx, bonder, msg.Coin)
	if err != nil {
		return nil, err
	}

	return &types.MsgBondInResponse{}, nil
}

func (k msgServer) Redeem(goCtx context.Context, msg *types.MsgRedeem) (*types.MsgRedeemResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	bonder, err := sdk.AccAddressFromBech32(msg.Bonder)
	if err != nil {
		return nil, err
	}

	err = k.keeper.RedeemDebt(ctx, bonder)
	if err != nil {
		return nil, err
	}

	return &types.MsgRedeemResponse{}, nil
}
