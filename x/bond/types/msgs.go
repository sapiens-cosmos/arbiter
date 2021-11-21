package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgBondIn = "bond_in"
)

var (
	_ sdk.Msg = &MsgBondIn{}
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
