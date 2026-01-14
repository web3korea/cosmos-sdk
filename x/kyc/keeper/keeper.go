package keeper

import (
	"encoding/json"
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/cosmos/cosmos-sdk/x/kyc/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   sdk.StoreKey
	memKey     sdk.StoreKey
	paramstore paramtypes.Subspace

	accountKeeper types.AccountKeeper
}

// NewKeeper creates new instances of the kyc Keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
) *Keeper {
	// Set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		paramstore:    ps,
		accountKeeper: accountKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetKYCRaw returns the KYC record for a given address
func (k Keeper) GetKYCRaw(ctx sdk.Context, address string) (types.KYC, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetKYCKey(address)
	bz := store.Get(key)
	if bz == nil {
		return types.KYC{}, false
	}

	var kyc types.KYC
	k.cdc.MustUnmarshalJSON(bz, &kyc)
	return kyc, true
}

// GetKYC returns the KYC record for a given address
func (k Keeper) GetKYC(ctx sdk.Context, address string) (types.KYC, error) {
	kyc, found := k.GetKYCRaw(ctx, address)
	if !found {
		return types.KYC{}, types.ErrKYCNotFound
	}
	return kyc, nil
}

// SetKYC sets the KYC record for a given address
func (k Keeper) SetKYC(ctx sdk.Context, kyc types.KYC) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetKYCKey(kyc.Address)
	bz := k.cdc.MustMarshalJSON(&kyc)
	store.Set(key, bz)
}

// DeleteKYC deletes the KYC record for a given address
func (k Keeper) DeleteKYC(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetKYCKey(address)
	store.Delete(key)
}

// HasKYC checks if a KYC record exists for a given address
func (k Keeper) HasKYC(ctx sdk.Context, address string) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetKYCKey(address)
	return store.Has(key)
}

// GetAllKYC returns all KYC records
func (k Keeper) GetAllKYC(ctx sdk.Context) []types.KYC {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KYCStoreKeyPrefix)
	defer iterator.Close()

	var kycRecords []types.KYC
	for ; iterator.Valid(); iterator.Next() {
		var kyc types.KYC
		k.cdc.MustUnmarshal(iterator.Value(), &kyc)
		kycRecords = append(kycRecords, kyc)
	}
	return kycRecords
}

// RegisterKYC registers a new KYC record
func (k Keeper) RegisterKYC(ctx sdk.Context, msg *types.MsgRegisterKYC) error {
	// Check if KYC already exists
	if k.HasKYC(ctx, msg.Sender) {
		return types.ErrKYCAlreadyExists
	}

	// Create new KYC record
	kyc := types.NewKYC(
		msg.Sender,
		msg.FullName,
		msg.DateOfBirth,
		msg.Country,
		msg.AddressInfo,
		msg.IDType,
		msg.IDNumber,
	)

	// Set KYC record
	k.SetKYC(ctx, kyc)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeKYCRegistered,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyStatus, kyc.Status),
		),
	)

	return nil
}

// ApproveKYC approves a KYC record
func (k Keeper) ApproveKYC(ctx sdk.Context, msg *types.MsgApproveKYC) error {
	// Check if sender is a validator
	if !k.IsValidator(ctx, msg.Sender) {
		return types.ErrUnauthorizedValidator
	}

	// Get KYC record
	kyc, found := k.GetKYC(ctx, msg.User)
	if !found {
		return types.ErrKYCNotFound
	}

	// Check if already approved
	if kyc.Status == types.StatusApproved {
		return types.ErrKYCAlreadyApproved
	}

	// Check if already rejected
	if kyc.Status == types.StatusRejected {
		return types.ErrKYCAlreadyRejected
	}

	// Update KYC status
	kyc.UpdateStatus(types.StatusApproved, msg.Sender, msg.Comments)
	k.SetKYC(ctx, kyc)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeKYCApproved,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.User),
			sdk.NewAttribute(types.AttributeKeyValidator, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyStatus, kyc.Status),
		),
	)

	return nil
}

// RejectKYC rejects a KYC record
func (k Keeper) RejectKYC(ctx sdk.Context, msg *types.MsgRejectKYC) error {
	// Check if sender is a validator
	if !k.IsValidator(ctx, msg.Sender) {
		return types.ErrUnauthorizedValidator
	}

	// Get KYC record
	kyc, found := k.GetKYC(ctx, msg.User)
	if !found {
		return types.ErrKYCNotFound
	}

	// Check if already approved
	if kyc.Status == types.StatusApproved {
		return types.ErrKYCAlreadyApproved
	}

	// Check if already rejected
	if kyc.Status == types.StatusRejected {
		return types.ErrKYCAlreadyRejected
	}

	// Update KYC status
	kyc.UpdateStatus(types.StatusRejected, msg.Sender, msg.Comments)
	k.SetKYC(ctx, kyc)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeKYCRejected,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.User),
			sdk.NewAttribute(types.AttributeKeyValidator, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyStatus, kyc.Status),
		),
	)

	return nil
}

// Validator management functions

// GetValidatorRaw returns the validator for a given address
func (k Keeper) GetValidatorRaw(ctx sdk.Context, address string) (types.Validator, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetValidatorKey(address)
	bz := store.Get(key)
	if bz == nil {
		return types.Validator{}, false
	}

	var validator types.Validator
	k.cdc.MustUnmarshal(bz, &validator)
	return validator, true
}

// GetValidator returns the validator for a given address
func (k Keeper) GetValidator(ctx sdk.Context, address string) (types.Validator, error) {
	validator, found := k.GetValidatorRaw(ctx, address)
	if !found {
		return types.Validator{}, types.ErrValidatorNotFound
	}
	return validator, nil
}

// SetValidator sets the validator for a given address
func (k Keeper) SetValidator(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetValidatorKey(validator.Address)
	bz := k.cdc.MustMarshal(&validator)
	store.Set(key, bz)
}

// DeleteValidator deletes the validator for a given address
func (k Keeper) DeleteValidator(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetValidatorKey(address)
	store.Delete(key)
}

// IsValidator checks if an address is a validator
func (k Keeper) IsValidator(ctx sdk.Context, address string) bool {
	validator, found := k.GetValidator(ctx, address)
	return found && validator.IsActive
}

// GetAllValidators returns all validators
func (k Keeper) GetAllValidators(ctx sdk.Context) []types.Validator {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorStoreKeyPrefix)
	defer iterator.Close()

	var validators []types.Validator
	for ; iterator.Valid(); iterator.Next() {
		var validator types.Validator
		k.cdc.MustUnmarshal(iterator.Value(), &validator)
		validators = append(validators, validator)
	}
	return validators
}

// AddValidator adds a new validator
func (k Keeper) AddValidator(ctx sdk.Context, address, name string, permissions []string) error {
	// Check if validator already exists
	if k.IsValidator(ctx, address) {
		return types.ErrValidatorAlreadyExists
	}

	validator := types.Validator{
		Address:     address,
		Name:        name,
		Permissions: permissions,
		IsActive:    true,
	}

	k.SetValidator(ctx, validator)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeValidatorAdded,
			sdk.NewAttribute(types.AttributeKeyAddress, address),
			sdk.NewAttribute(types.AttributeKeyName, name),
		),
	)

	return nil
}

// RemoveValidator removes a validator
func (k Keeper) RemoveValidator(ctx sdk.Context, address string) error {
	// Check if validator exists
	if !k.IsValidator(ctx, address) {
		return types.ErrValidatorNotFound
	}

	k.DeleteValidator(ctx, address)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeValidatorRemoved,
			sdk.NewAttribute(types.AttributeKeyAddress, address),
		),
	)

	return nil
}

// Parameter management functions

// GetParams gets the parameters for the kyc module
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the parameters for the kyc module
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// Genesis functions

// InitGenesis initializes the kyc module's state from a provided genesis state
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {
	// Set parameters
	k.SetParams(ctx, genState.Params)

	// Set KYC records
	for _, kyc := range genState.KYCRecords {
		k.SetKYC(ctx, kyc)
	}

	// Set validators
	for _, validator := range genState.Validators {
		k.SetValidator(ctx, validator)
	}
}

// ExportGenesis returns the kyc module's exported genesis
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := k.GetParams(ctx)
	kycRecords := k.GetAllKYC(ctx)
	validators := k.GetAllValidators(ctx)

	return types.NewGenesisState(params, kycRecords, validators)
}