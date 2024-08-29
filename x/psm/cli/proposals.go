package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cobra"

	"github.com/onomyprotocol/onomy/x/psm/types"
)

func NewCmdSubmitAddStableCoinProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-stable-coin [title] [description] [denom] [limit-total] [price] [fee_in] [fee_out] [proposer] [deposit]",
		Args:  cobra.ExactArgs(9),
		Short: "Submit an add stable coin proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			limitTotal, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return fmt.Errorf("value %s cannot constructs Int from string", args[3])
			}

			price := sdk.MustNewDecFromStr(args[4])
			feeIn := sdk.MustNewDecFromStr(args[5])
			feeOut := sdk.MustNewDecFromStr(args[6])
			from := sdk.MustAccAddressFromBech32(args[7])
			fmt.Println("kkkkk:from:", from.String())
			content := types.NewAddStableCoinProposal(
				args[0], args[1], args[2], limitTotal, price, feeIn, feeOut,
			)

			deposit, err := sdk.ParseCoinsNormalized(args[8])
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(&content, deposit, from)
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	// cmd.Flags().String(flags.FlagFrom, "", "Name or address of private key with which to sign")
	// flags.AddTxFlagsToCmd(cmd)

	return cmd
}
func NewCmdUpdatesStableCoinProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-limit-total-stable-coin [title] [description] [denom] [limit-total-update] [price] [fee_in] [fee_out] [deposit]",
		Args:  cobra.ExactArgs(8),
		Short: "Submit update limit total stable coin proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			limitTotalUpdate, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return fmt.Errorf("value %s cannot constructs Int from string", args[3])
			}
			price := sdk.MustNewDecFromStr(args[4])
			feeIn := sdk.MustNewDecFromStr(args[5])
			feeOut := sdk.MustNewDecFromStr(args[6])
			from := clientCtx.GetFromAddress()
			content := types.NewUpdatesStableCoinProposal(
				args[0], args[1], args[2], limitTotalUpdate, price, feeIn, feeOut,
			)

			deposit, err := sdk.ParseCoinsNormalized(args[7])
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(&content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}
