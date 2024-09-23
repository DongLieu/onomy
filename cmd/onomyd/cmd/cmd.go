// Package cmd contains cli wrapper for the onomy cli.
package cmd

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/server"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cosmoscmd"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/onomyprotocol/onomy/app"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibcprovidertypes "github.com/cosmos/interchain-security/x/ccv/provider/types"
)

// NewRootCmd initiates the cli for onomy chain.
func NewRootCmd() (*cobra.Command, cosmoscmd.EncodingConfig) {
	rootCmd, encodingConfig := cosmoscmd.NewRootCmd(
		app.Name,
		app.AccountAddressPrefix,
		app.DefaultNodeHome,
		app.Name,
		app.ModuleBasics,
		app.New,
	)
	// pull request #171 refactor: Remove ics. So we need re-register proto can read state
	RegisterInterfacesProvider(encodingConfig.InterfaceRegistry)

	rootCmd.AddCommand(
		server.RosettaCommand(encodingConfig.InterfaceRegistry, encodingConfig.Marshaler),
	)

	initRootCmd(rootCmd)
	return rootCmd, encodingConfig
}

// WrapBridgeCommands registers a sub-tree of gravity commands.
func WrapBridgeCommands(defaultNodeHome, rootCmd string, cmds []*cobra.Command) *cobra.Command {
	//nolint: exhaustivestruct
	cmd := &cobra.Command{
		Use:   rootCmd,
		Short: fmt.Sprintf("Manage %s bridge.", rootCmd),
		Long:  fmt.Sprintf("Manage %s bridge.", rootCmd),
	}

	for _, childCmd := range cmds {
		cmd.AddCommand(childCmd)
	}

	cmd.PersistentFlags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.PersistentFlags().String(flags.FlagKeyringDir, "", "The client Keyring directory; if omitted, the default 'home' directory will be used")
	cmd.PersistentFlags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.PersistentFlags().String(cli.OutputFlag, "text", "Output format (text|json)")

	return cmd
}

func initRootCmd(rootCmd *cobra.Command) {
	rootCmd.AddCommand(
		NewTestnetCmd(addModuleInitFlags),
	)
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
}

// // pull request #171 refactor: Remove ics. So we need re-register proto can read state
func RegisterInterfacesProvider(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&ibcprovidertypes.ConsumerAdditionProposal{},
	)
}
