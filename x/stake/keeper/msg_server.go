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
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	return &types.MsgJoinStakeResponse{}, nil
}

// func (server msgServer) CreateGauge(goCtx context.Context, msg *types.MsgCreateGauge) (*types.MsgCreateGaugeResponse, error) {
// ctx := sdk.UnwrapSDKContext(goCtx)
// owner, err := sdk.AccAddressFromBech32(msg.Owner)
// if err != nil {
// 	return nil, err
// }

// 	gaugeID, err := server.keeper.CreateGauge(ctx, msg.IsPerpetual, owner, msg.Coins, msg.DistributeTo, msg.StartTime, msg.NumEpochsPaidOver)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
// 	}

// 	ctx.EventManager().EmitEvents(sdk.Events{
// 		sdk.NewEvent(
// 			types.TypeEvtCreateGauge,
// 			sdk.NewAttribute(types.AttributeGaugeID, utils.Uint64ToString(gaugeID)),
// 		),
// 	})

// 	return &types.MsgCreateGaugeResponse{}, nil
// }
