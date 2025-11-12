package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/onomyprotocol/onomy/x/gnark/types"
)

const (
	flagVerifyingKey = "verifying-key"
	flagCircuitFile  = "circuit"
	flagProofFile    = "proof"
	flagWitnessFile  = "public-witness"
)

// TxCmd bundles all gnark transaction subcommands under a single root command.
func TxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Gnark transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewCreateCircuitCmd(),
		NewVerifyProofCmd(),
	)

	return cmd
}

// NewCreateCircuitCmd wires the MsgCreateCircuit flow into the CLI.
func NewCreateCircuitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-circuit [circuit-id] [curve-id]",
		Short: "Upload a gnark verifying key and circuit",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			vkPath, err := cmd.Flags().GetString(flagVerifyingKey)
			if err != nil {
				return err
			}
			vkBytes, err := os.ReadFile(vkPath)
			if err != nil {
				return fmt.Errorf("read verifying key: %w", err)
			}

			circuitPath, err := cmd.Flags().GetString(flagCircuitFile)
			if err != nil {
				return err
			}
			circuitBytes, err := os.ReadFile(circuitPath)
			if err != nil {
				return fmt.Errorf("read circuit: %w", err)
			}

			msg := types.NewMsgCreateCircuit(
				clientCtx.GetFromAddress().String(),
				args[0],
				args[1],
				vkBytes,
				circuitBytes,
			)

			if err := tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().String(flagVerifyingKey, "", "path to the verifying key file")
	cmd.Flags().String(flagCircuitFile, "", "path to the circuit (.r1cs) file")
	_ = cmd.MarkFlagRequired(flagVerifyingKey)
	_ = cmd.MarkFlagRequired(flagCircuitFile)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewVerifyProofCmd wires the MsgVerifyProof flow into the CLI.
func NewVerifyProofCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-proof [circuit-id]",
		Short: "Submit a proof for verification against a stored circuit",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proofPath, err := cmd.Flags().GetString(flagProofFile)
			if err != nil {
				return err
			}
			proofBytes, err := os.ReadFile(proofPath)
			if err != nil {
				return fmt.Errorf("read proof: %w", err)
			}

			witnessPath, err := cmd.Flags().GetString(flagWitnessFile)
			if err != nil {
				return err
			}
			witnessBytes, err := os.ReadFile(witnessPath)
			if err != nil {
				return fmt.Errorf("read public witness: %w", err)
			}

			msg := types.NewMsgVerifyProof(
				clientCtx.GetFromAddress().String(),
				args[0],
				proofBytes,
				witnessBytes,
			)

			if err := tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().String(flagProofFile, "", "path to the Groth16 proof file")
	cmd.Flags().String(flagWitnessFile, "", "path to the public witness file")
	_ = cmd.MarkFlagRequired(flagProofFile)
	_ = cmd.MarkFlagRequired(flagWitnessFile)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
