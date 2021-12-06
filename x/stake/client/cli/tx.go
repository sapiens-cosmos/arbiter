package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sapiens-cosmos/arbiter/x/stake/types"
)

func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Generalized automated market maker transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewStakeCmd(),
		NewClaimCommand(),
	)
	return txCmd
}

func NewStakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stake [token-in]",
		Short: "stake designated amount",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			stakeAmount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			msg := types.NewMsgJoinStake(
				clientCtx.GetFromAddress(),
				stakeAmount,
			)

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewClaimCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim [token-in]",
		Short: "claim designated amount",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			claimAmount, ok := sdk.NewIntFromString(args[0])
			if !ok {
				return fmt.Errorf("invalid claim amount")
			}

			msg := types.NewMsgClaim(
				clientCtx.GetFromAddress(),
				claimAmount,
			)

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
