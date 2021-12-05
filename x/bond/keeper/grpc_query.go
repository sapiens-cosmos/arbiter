package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sapiens-cosmos/arbiter/x/bond/types"
)

type queryServer struct {
	keeper Keeper
}

// NewQueryServerImpl returns an implementation of the bond QueryServer interface
// for the provided Keeper.
func NewQueryServerImpl(keeper Keeper) types.QueryServer {
	return &queryServer{keeper: keeper}
}

func (q queryServer) Redeemable(ctx context.Context, req *types.QueryRedeemableRequest) (*types.QueryRedeemabeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Bonder == "" {
		return nil, status.Error(codes.InvalidArgument, "bonder cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	address, err := sdk.AccAddressFromBech32(req.Bonder)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %s", err.Error())
	}

	redeemable, err := q.keeper.RedeeambleDebt(sdkCtx, address)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRedeemabeResponse{Coin: &redeemable}, nil
}

func (q queryServer) RiskFreePrice(ctx context.Context, req *types.QueryRiskFreePriceRequest) (*types.QueryRiskFreePriceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.BondDenom == "" {
		return nil, status.Error(codes.InvalidArgument, "bond denom cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	rfp, err := q.keeper.GetRiskFreePrice(sdkCtx, req.BondDenom)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRiskFreePriceResponse{Price: &rfp}, nil
}

func (q queryServer) Premium(ctx context.Context, req *types.QueryPremiumRequest) (*types.QueryPremiumResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.BondDenom == "" {
		return nil, status.Error(codes.InvalidArgument, "bond denom cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	premium, err := q.keeper.GetPremium(sdkCtx, req.BondDenom)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPremiumResponse{Premium: &premium}, nil
}

func (q queryServer) BondInfo(ctx context.Context, req *types.QueryBondInfoRequest) (*types.QueryBondInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.BondDenom == "" {
		return nil, status.Error(codes.InvalidArgument, "bond denom cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	rfp, err := q.keeper.GetRiskFreePrice(sdkCtx, req.BondDenom)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	premium, err := q.keeper.GetPremium(sdkCtx, req.BondDenom)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Executing Price = RiskFreePrice * Premium {Premium ≥ 1}
	executingPrice := rfp.Mul(premium)

	return &types.QueryBondInfoResponse{
		RiskFreePrice:  &rfp,
		Premium:        &premium,
		ExecutingPrice: &executingPrice,
	}, nil
}

func (q queryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	params := q.keeper.GetParams(sdkCtx)

	return &types.QueryParamsResponse{Params: params}, nil
}
