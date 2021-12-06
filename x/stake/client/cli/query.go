package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/sapiens-cosmos/arbiter/x/stake/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group epochs queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetCmdBalance(),
		GetCmdStaked(),
		GetCmdTimeUnitlRebase(),
	)

	return cmd
}

func GetCmdBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance <address>",
		Short: "Query balance of address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query balance of address.

Example:
$ %s query stake balance <address>
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Balance(cmd.Context(), &types.QueryBalanceRequest{Sender: args[0]})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdStaked() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "staked <address>",
		Short: "Query stake of address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query staked of address.

Example:
$ %s query stake staked <address>
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Staked(cmd.Context(), &types.QueryStakedRequest{Sender: args[0]})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdTimeUnitlRebase() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "time-until-rebase",
		Short: "Returns block left until rebase",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query block unitl rebase.

Example:
$ %s query stake time-until-rebase
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.TimeUntilRebase(cmd.Context(), &types.QueryTimeUntilRebaseRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
