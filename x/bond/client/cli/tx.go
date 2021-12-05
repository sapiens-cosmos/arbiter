package cli

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/sapiens-cosmos/arbiter/x/bond/types"
)

// GetTxCmd returns a root CLI command handler for all x/bond transaction commands.
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Bond transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(NewBondInCmd())
	txCmd.AddCommand(NewRedeemCmd())

	return txCmd
}

// NewBondInCmd returns a CLI command handler for creating a MsgBondIn transaction.
func NewBondInCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "bond-in [from_key_or_address] [coin]",
		Short: `Bond asset to receive discounted base coin. Note, the'--from' flag is
ignored as it is implied from [from_key_or_address].`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.Flags().Set(flags.FlagFrom, args[0])
			if err != nil {
				panic(err)
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coin, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgBondIn(clientCtx.GetFromAddress(), coin)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewRedeemCmd returns a CLI command handler for creating a MsgRedeem transaction.
func NewRedeemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "redeem [from_key_or_address]",
		Short: `Redeem debt from bond. Note, the'--from' flag is
ignored as it is implied from [from_key_or_address].`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.Flags().Set(flags.FlagFrom, args[0])
			if err != nil {
				panic(err)
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRedeem(clientCtx.GetFromAddress())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
