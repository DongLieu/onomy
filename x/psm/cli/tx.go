package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/onomyprotocol/onomy/x/psm/types"
	"github.com/spf13/cobra"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2, // nolint:gomnd
		RunE:                       client.ValidateCmd,
	}

	// cmd.AddCommand(NewCmdUpdatesStableCoinProposal())
	// cmd.AddCommand(NewCmdSubmitAddStableCoinProposal())
	cmd.AddCommand(NewSwapToStablecoinCmd())
	cmd.AddCommand(NewSwapToISTCmd())

	return cmd
}

func NewSwapToISTCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap-to-ist [stablecoin]",
		Args:  cobra.ExactArgs(1),
		Short: "swap stablecoin to $ist ",
		Long: `swap stablecoin to $ist.

			Example:
			$ onomyd tx psm swap-to-ist 1000usdt --from mykey
	`,

		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			stablecoin, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			addr := clientCtx.GetFromAddress()
			msg := types.NewMsgSwapToIST(addr.String(), &stablecoin)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewSwapToStablecoinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap-to-stablecoin [stable-coin-type] [amount-ist]",
		Args:  cobra.ExactArgs(2),
		Short: "swap $ist to stablecoin ",
		Long: `swap $ist to stablecoin.

			Example:
			$ onomyd tx psm swap-to-stablecoin usdt 10000ist --from mykey
	`,

		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			ISTcoin, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			addr := clientCtx.GetFromAddress()
			msg := types.NewMsgSwapToStablecoin(addr.String(), args[0], ISTcoin.Amount)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
