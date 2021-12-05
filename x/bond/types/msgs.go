package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgBondIn = "bond_in"
	TypeMsgRedeem = "redeem"
)

var (
	_ sdk.Msg = &MsgBondIn{}
	_ sdk.Msg = &MsgRedeem{}
)

func NewMsgBondIn(bonder sdk.AccAddress, coin sdk.Coin) *MsgBondIn {
	return &MsgBondIn{
		Bonder: bonder.String(),
		Coin:   coin,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgBondIn) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgBondIn) Type() string { return TypeMsgBondIn }

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
func (msg MsgBondIn) GetSigners() []sdk.AccAddress {
	bonder, err := sdk.AccAddressFromBech32(msg.Bonder)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{bonder}
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgBondIn) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgBondIn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Bonder)
	if err != nil {
		return err
	}

	return msg.Coin.Validate()
}

func NewMsgRedeem(bonder sdk.AccAddress) *MsgRedeem {
	return &MsgRedeem{
		Bonder: bonder.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgRedeem) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgRedeem) Type() string { return TypeMsgRedeem }

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
func (msg MsgRedeem) GetSigners() []sdk.AccAddress {
	bonder, err := sdk.AccAddressFromBech32(msg.Bonder)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{bonder}
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgRedeem) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgRedeem) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Bonder)
	if err != nil {
		return err
	}

	return nil
}
