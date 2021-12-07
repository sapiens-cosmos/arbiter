package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sapiens-cosmos/arbiter/x/stake/types"
)

type queryServer struct {
	keeper Keeper
}

// NewQueryServerImpl returns an implementation of the bond QueryServer interface
// for the provided Keeper.
func NewQueryServerImpl(keeper Keeper) types.QueryServer {
	return &queryServer{keeper: keeper}
}

func (q queryServer) Balance(ctx context.Context, req *types.QueryBalanceRequest) (*types.QueryBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Sender == "" {
		return nil, status.Error(codes.InvalidArgument, "sender cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	address, err := sdk.AccAddressFromBech32(req.Sender)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %s", err.Error())
	}

	balance := q.keeper.GetBalance(sdkCtx, address)

	return &types.QueryBalanceResponse{Balance: &balance}, nil
}

func (q queryServer) Staked(ctx context.Context, req *types.QueryStakedRequest) (*types.QueryStakedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Sender == "" {
		return nil, status.Error(codes.InvalidArgument, "sender cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	address, err := sdk.AccAddressFromBech32(req.Sender)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %s", err.Error())
	}

	staked, err := q.keeper.GetStakedTokenByAddress(sdkCtx, address)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryStakedResponse{Staked: &staked}, nil
}

func (q queryServer) TimeUntilRebase(ctx context.Context, req *types.QueryTimeUntilRebaseRequest) (*types.QueryTimeUntilRebaseResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockUntilRebase := q.keeper.GetBlockUntilRebase(sdkCtx)

	return &types.QueryTimeUntilRebaseResponse{BlockUntilRebase: blockUntilRebase}, nil
}

func (q queryServer) RewardYield(ctx context.Context, req *types.QueryRewardYieldRequest) (*types.QueryRewardYieldResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	rewardYield := q.keeper.GetRewardYield(sdkCtx)

	return &types.QueryRewardYieldResponse{RewardYield: &rewardYield}, nil
}

func (q queryServer) StakeInfo(ctx context.Context, req *types.QueryStakeInfoRequest) (*types.QueryStakeInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Sender == "" {
		return nil, status.Error(codes.InvalidArgument, "sender cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	address, err := sdk.AccAddressFromBech32(req.Sender)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %s", err.Error())
	}

	balance := q.keeper.GetBalance(sdkCtx, address)
	staked, err := q.keeper.GetStakedTokenByAddress(sdkCtx, address)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	blockUntilRebase := q.keeper.GetBlockUntilRebase(sdkCtx)
	rewardYield := q.keeper.GetRewardYield(sdkCtx)

	totalStaked := q.keeper.GetTotalStaked(sdkCtx)

	return &types.QueryStakeInfoResponse{Balance: &balance, Staked: &staked, BlockUntilRebase: blockUntilRebase, RewardYield: &rewardYield, TotalStaked: &totalStaked}, nil
}

func (q queryServer) TotalReserve(ctx context.Context, req *types.QueryTotalReserveRequest) (*types.QueryTotalReserveResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	totalReserve := q.keeper.GetTotalReserve(sdkCtx).Int64()

	return &types.QueryTotalReserveResponse{TotalReserve: totalReserve}, nil
}

func (q queryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	params := q.keeper.GetParams(sdkCtx)

	return &types.QueryParamsResponse{Params: params}, nil
}
