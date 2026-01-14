package cli

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/cosmos/cosmos-sdk/x/kyc/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	kycTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "KYC transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	kycTxCmd.AddCommand(
		NewRegisterKYCCmd(),
		NewApproveKYCCmd(),
		NewRejectKYCCmd(),
	)

	return kycTxCmd
}

// NewRegisterKYCCmd implements the register-kyc command
func NewRegisterKYCCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-kyc [full-name] [date-of-birth] [country] [address-info] [id-type] [id-number]",
		Short: "Register a new KYC record",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			dateOfBirth, err := time.Parse("2006-01-02", args[1])
			if err != nil {
				return err
			}

			msg := &types.MsgRegisterKYC{
				Sender:      clientCtx.GetFromAddress().String(),
				FullName:    args[0],
				DateOfBirth: dateOfBirth,
				Country:     args[2],
				AddressInfo: args[3],
				IDType:      args[4],
				IDNumber:    args[5],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewApproveKYCCmd implements the approve-kyc command
func NewApproveKYCCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve-kyc [user-address] [comments]",
		Short: "Approve a KYC record",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgApproveKYC{
				Sender:   clientCtx.GetFromAddress().String(),
				User:     args[0],
				Comments: args[1],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewRejectKYCCmd implements the reject-kyc command
func NewRejectKYCCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reject-kyc [user-address] [comments]",
		Short: "Reject a KYC record",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgRejectKYC{
				Sender:   clientCtx.GetFromAddress().String(),
				User:     args[0],
				Comments: args[1],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}