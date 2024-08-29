package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// ParamSubspace defines the expected Subspace interface.
type ParamSubspace interface {
	HasKeyTable() bool
	WithKeyTable(paramtypes.KeyTable) paramtypes.Subspace
	Get(sdk.Context, []byte, interface{})
	GetParamSet(sdk.Context, paramtypes.ParamSet)
	SetParamSet(sdk.Context, paramtypes.ParamSet)
}

// AccountKeeper defines the contract required for account APIs.
type AccountKeeper interface {
	GetModuleAddress(string) sdk.AccAddress
}

// BankKeeper defines the contract needed to be fulfilled for banking and supply dependencies.
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	GetAllBalances(sdk.Context, sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoins(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	LockedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// DistributionKeeper expected distribution keeper.
type DistributionKeeper interface {
	HasDelegatorStartingInfo(sdk.Context, sdk.ValAddress, sdk.AccAddress) bool
	WithdrawDelegationRewards(sdk.Context, sdk.AccAddress, sdk.ValAddress) (sdk.Coins, error)
}

// GovKeeper expected gov keeper.
type GovKeeper interface {
	AddVote(sdk.Context, uint64, sdk.AccAddress, govtypes.WeightedVoteOptions) error
	GetVote(sdk.Context, uint64, sdk.AccAddress) (govtypes.Vote, bool)
	IterateProposals(sdk.Context, func(proposal govtypes.Proposal) bool)
}

// MintKeeper expected mint keeper.
type MintKeeper interface {
	GetMinter(ctx sdk.Context) (minter minttypes.Minter)
	GetParams(ctx sdk.Context) (params minttypes.Params)
}

// StakingKeeper expected staking keeper.
type StakingKeeper interface {
	BondDenom(sdk.Context) string
	Delegate(sdk.Context, sdk.AccAddress, sdk.Int, stakingtypes.BondStatus, stakingtypes.Validator, bool) (sdk.Dec, error)
	GetDelegation(sdk.Context, sdk.AccAddress, sdk.ValAddress) (stakingtypes.Delegation, bool)
	GetAllValidators(sdk.Context) []stakingtypes.Validator
	UnbondAndUndelegateCoins(sdk.Context, sdk.AccAddress, sdk.ValAddress, sdk.Dec) (sdk.Int, error)
}

type PriceFeeds interface {
}

type OracleKeeper interface {
	GetMultiplier(ctx sdk.Context, denom string) sdk.Dec
}
