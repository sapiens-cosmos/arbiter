package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgJoinStake = "stake"
)

var _ sdk.Msg = &MsgJoinStake{}

func NewMsgJoinStake(owner sdk.AccAddress, tokenIn sdk.Coin) *MsgJoinStake {
	return &MsgJoinStake{
		Sender:  owner.String(),
		TokenIn: tokenIn,
	}
}

func (m MsgJoinStake) Route() string { return RouterKey }
func (m MsgJoinStake) Type() string  { return TypeMsgJoinStake }
func (m MsgJoinStake) ValidateBasic() error {
	if m.Sender == "" {
		return errors.New("owner should be set")
	}
	return nil
}
func (m MsgJoinStake) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgJoinStake) GetSigners() []sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{owner}
}

func NewMsgClaim(owner sdk.AccAddress, tokenIn sdk.Int) *MsgClaim {
	return &MsgClaim{
		Sender:  owner.String(),
		TokenIn: tokenIn,
	}
}

func (m MsgClaim) Route() string { return RouterKey }
func (m MsgClaim) Type() string  { return TypeMsgJoinStake }
func (m MsgClaim) ValidateBasic() error {
	if m.Sender == "" {
		return errors.New("owner should be set")
	}
	return nil
}
func (m MsgClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgClaim) GetSigners() []sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{owner}
}
