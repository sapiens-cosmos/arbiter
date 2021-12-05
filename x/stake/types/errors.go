package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrNoLock = sdkerrors.Register(ModuleName, 1, "no lock registered for address")
)
