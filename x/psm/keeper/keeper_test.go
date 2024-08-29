package keeper_test

import (
	"log"
	"os/exec"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/onomyprotocol/onomy/testutil/simapp"

	"github.com/onomyprotocol/onomy/x/psm/types"
	"github.com/stretchr/testify/require"
)

var (
	limitTotal      = sdk.NewInt(1000000000000000000)
	fee             = sdk.MustNewDecFromStr("0.01")
	price           = sdk.MustNewDecFromStr("1")
	priceAcceptable = sdk.MustNewDecFromStr("0.99")
)

func TestStablecoin(t *testing.T) {
	app := simapp.Setup()
	k := app.OnomyApp().PSMKeeper
	ctx := app.NewContext()

	s := []types.Stablecoin{
		{
			Denom:      "usdt",
			LimitTotal: limitTotal,
			Price:      price,
			FeeIn:      fee,
			FeeOut:     fee,
		},
		{
			Denom:      "usdc",
			LimitTotal: limitTotal,
			Price:      priceAcceptable,
			FeeIn:      fee,
			FeeOut:     fee,
		},
	}

	k.SetStablecoin(ctx, s[0])
	k.SetStablecoin(ctx, s[1])

	s0, f := k.GetStablecoin(ctx, "usdt")
	require.True(t, f)
	require.Equal(t, s0.Denom, "usdt")
	require.Equal(t, s0.LimitTotal, limitTotal)
	require.Equal(t, s0.Price, price)

	s1, f := k.GetStablecoin(ctx, "usdc")
	require.True(t, f)
	require.Equal(t, s1.Denom, "usdc")
	require.Equal(t, s1.LimitTotal, limitTotal)
	require.Equal(t, s1.Price, priceAcceptable)

	var count = 0
	k.IterateStablecoin(ctx, func(red types.Stablecoin) (stop bool) {
		count += 1
		return false
	})
	require.Equal(t, count, 2)

	cmd := exec.Command("sh", "-c", "testchaind start")

	// Thực thi lệnh và kiểm tra lỗi
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to execute command: %s", err)
	}

	log.Println("Command executed successfully")

	time.Sleep(time.Second * 20)

	cmd = exec.Command("sh", "-c", "killall testchaind || true")
	// Thực thi lệnh và kiểm tra lỗi
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to execute commandd1: %s", err)
	}
}
