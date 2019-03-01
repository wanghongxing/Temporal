package commands

import (
	"context"
	"fmt"
	"path"
	"strconv"

	cli "gx/ipfs/Qmc1AtgBdoUHP8oYSqU81NRYdzohmF45t5XNwVMvhCxsBA/cli"

	"gx/ipfs/QmckeQ2zrYLAXoSHYTGn5BDdb22BqbUoHEHm8KZ9YWRxd1/iptb/testbed"
)

var ShellCmd = cli.Command{
	Category:  "CORE",
	Name:      "shell",
	Usage:     "starts a shell within the context of node",
	ArgsUsage: "<node>",
	Action: func(c *cli.Context) error {
		flagRoot := c.GlobalString("IPTB_ROOT")
		flagTestbed := c.GlobalString("testbed")

		if !c.Args().Present() {
			return NewUsageError("shell takes exactly 1 argument")
		}

		i, err := strconv.Atoi(c.Args().First())
		if err != nil {
			return fmt.Errorf("parse err: %s", err)
		}

		tb := testbed.NewTestbed(path.Join(flagRoot, "testbeds", flagTestbed))

		nodes, err := tb.Nodes()
		if err != nil {
			return err
		}

		return nodes[i].Shell(context.Background(), nodes)
	},
}
