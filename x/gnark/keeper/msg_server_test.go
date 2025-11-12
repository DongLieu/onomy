package keeper_test

import (
	"os"
	"path/filepath"
	"testing"

	colog "cosmossdk.io/log"
	store "cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	_ "github.com/onomyprotocol/onomy/app"
	"github.com/onomyprotocol/onomy/x/gnark/keeper"
	"github.com/onomyprotocol/onomy/x/gnark/types"
)

func setupKeeper(t *testing.T) (keeper.Keeper, sdk.Context) {
	t.Helper()

	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, colog.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	types.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)

	k := keeper.NewKeeper(cdc, storeKey)
	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, colog.NewNopLogger())
	return k, ctx
}

const fixturesDir = "../../../store"

func testAddress() string {
	priv := secp256k1.GenPrivKey()
	return sdk.AccAddress(priv.PubKey().Address()).String()
}

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Clean(path))
	require.NoError(t, err)
	return data
}

func TestMsgServerCreateAndVerify(t *testing.T) {
	k, sdkCtx := setupKeeper(t)
	msgSrv := keeper.NewMsgServerImpl(k)

	creator := testAddress()

	createMsg := &types.MsgCreateCircuit{
		Creator:      creator,
		CircuitId:    "poly",
		CurveId:      "BN254",
		VerifyingKey: mustReadFile(t, filepath.Join(fixturesDir, "poly_verifying.key")),
		Circuit:      mustReadFile(t, filepath.Join(fixturesDir, "poly_circuit.r1cs")),
	}

	_, err := msgSrv.CreateCircuit(sdk.WrapSDKContext(sdkCtx), createMsg)
	require.NoError(t, err)

	circuit, found := k.GetCircuit(sdkCtx, "poly")
	require.True(t, found)
	require.Equal(t, uint64(0), circuit.Height)

	verifyMsg := &types.MsgVerifyProof{
		Creator:       creator,
		CircuitId:     "poly",
		Proof:         mustReadFile(t, filepath.Join(fixturesDir, "poly_proof.bin")),
		PublicWitness: mustReadFile(t, filepath.Join(fixturesDir, "poly_public.wtns")),
	}

	resp, err := msgSrv.VerifyProof(sdk.WrapSDKContext(sdkCtx), verifyMsg)
	require.NoError(t, err)
	require.True(t, resp.Valid)
}
