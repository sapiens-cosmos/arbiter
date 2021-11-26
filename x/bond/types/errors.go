package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrNoBondPolicy = sdkerrors.Register(ModuleName, 1, "there is no bond policy")
	ErrNoDebt       = sdkerrors.Register(ModuleName, 2, "there is no debt")
)
