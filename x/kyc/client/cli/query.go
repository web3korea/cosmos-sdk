package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/cosmos/cosmos-sdk/x/kyc/types"
)

// GetQueryCmd returns the query commands for this module
func GetQueryCmd() *cobra.Command {
	kycQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the kyc module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	kycQueryCmd.AddCommand(
		GetCmdQueryKYC(),
		GetCmdListKYC(),
		GetCmdQueryValidator(),
		GetCmdListValidators(),
		GetCmdQueryParams(),
	)

	return kycQueryCmd
}

// GetCmdQueryKYC implements the query kyc command
func GetCmdQueryKYC() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kyc [address]",
		Short: "Query a KYC record by address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryGetKYCRequest{
				Address: args[0],
			}

			res, err := queryClient.GetKYC(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdListKYC implements the query list-kyc command
func GetCmdListKYC() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-kyc",
		Short: "Query all KYC records",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryListKYCRequest{}

			res, err := queryClient.ListKYC(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryValidator implements the query validator command
func GetCmdQueryValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator [address]",
		Short: "Query a validator by address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryGetValidatorRequest{
				Address: args[0],
			}

			res, err := queryClient.GetValidator(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdListValidators implements the query list-validators command
func GetCmdListValidators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-validators",
		Short: "Query all validators",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryListValidatorsRequest{}

			res, err := queryClient.ListValidators(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryParams implements the query params command
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current kyc parameters",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryParamsRequest{}

			res, err := queryClient.Params(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}