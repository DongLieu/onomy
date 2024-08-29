package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/onomyprotocol/onomy/x/psm/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd(_ string) *cobra.Command {
	// Group dao queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2, // nolint:gomnd
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdShowParams())
	cmd.AddCommand(CmdCheckStablecoin())

	return cmd
}

// CmdShowParams returns CmdShowParams cobra.Command.
func CmdShowParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-params",
		Short: "shows the parameters of the module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdShowParams returns CmdShowParams cobra.Command.
func CmdCheckStablecoin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stablecoin [denom]",
		Short: "shows info price, fee in, fee out, swapable quantity of the stablecoin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Stablecoin(context.Background(), &types.QueryStablecoinRequest{Denom: args[0]})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
