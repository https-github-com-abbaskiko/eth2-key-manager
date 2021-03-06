package account_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/util/printer"
)

func TestAccountList(t *testing.T) {
	t.Run("Successfully list accounts", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"list",
			"--storage=7b226163636f756e7473223a2237623764222c226174744d656d6f7279223a2237623764222c226e6574776f726b223a223664363136393665222c2270726f706f73616c4d656d6f7279223a2237623764222c2277616c6c6574223a2237623232363936343232336132323336363333383333333633353333333732643338333533363632326433343631363436333264363136313337333232643336363533393339333536363631333736363339333933393232326332323639366536343635373834643631373037303635373232323361376237643263323237343739373036353232336132323438343432323764227d",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})

	t.Run("Fail to JSON un-marshal", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"list",
			"--storage=7b226163636f756e7473223a2237623764222c226174744d656d6f7279223a2237623764222c2270726f706f",
		})
		err := cmd.RootCmd.Execute()
		require.Error(t, err)
		require.EqualError(t, err, "failed to JSON un-marshal storage: unexpected end of JSON input")
	})
}
