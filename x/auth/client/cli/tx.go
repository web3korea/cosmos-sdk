package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// NewTxCmd returns a root CLI command handler for all x/auth transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Auth transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewSetKycVerifiedTxCmd(),
	)

	return txCmd
}

// NewSetKycVerifiedTxCmd returns a CLI command handler for creating a MsgSetKycVerified transaction.
func NewSetKycVerifiedTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-kyc-verified [admin_key_or_address] [account_address] [kyc_verified]",
		Short: "Set KYC verification status for an account (admin only).",
		Long: `Set KYC verification status for an account. Only admin accounts can perform this action.
Note, the '--from' flag is ignored as it is implied from [admin_key_or_address].
When using '--dry-run' a key name cannot be used, only a bech32 address.
`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[0])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			adminAddr := clientCtx.GetFromAddress()
			accountAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			kycVerified := args[2] == "true"

			msg := types.NewMsgSetKycVerified(adminAddr, accountAddr, kycVerified)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}