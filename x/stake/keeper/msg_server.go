package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sapiens-cosmos/arbiter/x/stake/types"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServerImpl returns an instance of MsgServer
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

var _ types.MsgServer = msgServer{}

func (server msgServer) JoinStake(goCtx context.Context, msg *types.MsgJoinStake) (*types.MsgJoinStakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	server.keeper.JoinStake(ctx, msg.Sender, msg.TokenIn)

	return &types.MsgJoinStakeResponse{}, nil
}

func (server msgServer) Claim(goCtx context.Context, msg *types.MsgClaim) (*types.MsgClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	server.keeper.Claim(ctx, msg.Sender, msg.TokenIn)

	return &types.MsgClaimResponse{}, nil
}
