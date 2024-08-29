package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/onomyprotocol/onomy/x/psm/types"
)

type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      sdk.StoreKey
	memKey        sdk.StoreKey
	ps            types.ParamSubspace
	bankKeeper    types.BankKeeper
	OracleKeeper  types.OracleKeeper
	accountKeeper types.AccountKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps types.ParamSubspace,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	// ensure dao module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		ps:            ps,
		bankKeeper:    bankKeeper,
		accountKeeper: accountKeeper,
	}
}

// Logger returns keeper logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SetStablecoin(ctx sdk.Context, s types.Stablecoin) {
	store := ctx.KVStore(k.storeKey)

	key := types.GetKeyStableCoin(s.Denom)
	bz := k.cdc.MustMarshal(&s)

	store.Set(key, bz)
}

func (k Keeper) GetStablecoin(ctx sdk.Context, denom string) (types.Stablecoin, bool) {
	store := ctx.KVStore(k.storeKey)

	key := types.GetKeyStableCoin(denom)

	bz := store.Get(key)
	if bz == nil {
		return types.Stablecoin{}, false
	}

	var token types.Stablecoin
	k.cdc.MustUnmarshal(bz, &token)

	return token, true
}

func (k Keeper) IterateStablecoin(ctx sdk.Context, cb func(red types.Stablecoin) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.KeyStableCoin)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var token types.Stablecoin
		k.cdc.MustUnmarshal(iterator.Value(), &token)
		if cb(token) {
			break
		}
	}
}

func (k Keeper) SwaptoIST(ctx sdk.Context, addr sdk.AccAddress, stablecoin sdk.Coin) (sdk.Int, sdk.Coin, error) {
	asset := k.bankKeeper.GetBalance(ctx, addr, stablecoin.Denom)

	if asset.Amount.LT(stablecoin.Amount) {
		return sdk.ZeroInt(), sdk.Coin{}, fmt.Errorf("insufficient balance")
	}

	multiplier, err := k.GetPrice(ctx, stablecoin.Denom)
	if err != nil || multiplier.IsZero() {
		return sdk.Int{}, sdk.Coin{}, err
	}

	amountIST := multiplier.Mul(stablecoin.Amount.ToDec()).RoundInt()

	fee, err := k.PayFeesIn(ctx, amountIST, stablecoin.Denom)
	if err != nil {
		return sdk.Int{}, sdk.Coin{}, err
	}

	receiveAmountIST := amountIST.Sub(fee)
	return receiveAmountIST, sdk.NewCoin(types.InterStableToken, fee), nil
}

func (k Keeper) PayFeesIn(ctx sdk.Context, amount sdk.Int, denom string) (sdk.Int, error) {
	ratioSwapInFees, err := k.GetFeeIn(ctx, denom)
	if err != nil {
		return sdk.Int{}, err
	}
	fee := ratioSwapInFees.MulInt(amount).RoundInt()
	return fee, nil
}

func (k Keeper) SwapToStablecoin(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Int, toDenom string) (sdk.Int, sdk.Coin, error) {
	asset := k.bankKeeper.GetBalance(ctx, addr, types.InterStableToken)

	if asset.Amount.LT(amount) {
		return sdk.ZeroInt(), sdk.Coin{}, fmt.Errorf("insufficient balance")
	}

	multiplier, err := k.GetPrice(ctx, toDenom)
	if err != nil || multiplier.IsZero() {
		return sdk.Int{}, sdk.Coin{}, err
	}
	amountStablecoin := amount.ToDec().Quo(multiplier).RoundInt()

	fee, err := k.PayFeesOut(ctx, amountStablecoin, toDenom)
	if err != nil {
		return sdk.Int{}, sdk.Coin{}, err
	}

	receiveAmount := amountStablecoin.Sub(fee)
	return receiveAmount, sdk.NewCoin(toDenom, fee), nil
}

func (k Keeper) PayFeesOut(ctx sdk.Context, amount sdk.Int, denom string) (sdk.Int, error) {
	ratioSwapOutFees, err := k.GetFeeOut(ctx, denom)
	if err != nil {
		return sdk.Int{}, err
	}

	fee := ratioSwapOutFees.MulInt(amount).RoundInt()
	return fee, nil
}

func (k Keeper) GetTotalLimitWithDenomStablecoin(ctx sdk.Context, denom string) (sdk.Int, error) {
	s, found := k.GetStablecoin(ctx, denom)
	if !found {
		return sdk.Int{}, fmt.Errorf("not found Stable coin %s", denom)
	}
	return s.LimitTotal, nil
}

func (k Keeper) GetPrice(ctx sdk.Context, denom string) (sdk.Dec, error) {
	s, found := k.GetStablecoin(ctx, denom)
	if !found {
		return sdk.Dec{}, fmt.Errorf("not found Stable coin %s", denom)
	}
	return s.Price, nil
}

func (k Keeper) GetFeeIn(ctx sdk.Context, denom string) (sdk.Dec, error) {
	s, found := k.GetStablecoin(ctx, denom)
	if !found {
		return sdk.Dec{}, fmt.Errorf("not found Stable coin %s", denom)
	}
	return s.FeeIn, nil
}

func (k Keeper) GetFeeOut(ctx sdk.Context, denom string) (sdk.Dec, error) {
	s, found := k.GetStablecoin(ctx, denom)
	if !found {
		return sdk.Dec{}, fmt.Errorf("not found Stable coin %s", denom)
	}
	return s.FeeOut, nil
}
