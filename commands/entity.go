package commands

import (
	"github.com/codegangsta/cli"
	"github.com/brooklyncentral/brooklyn-cli/api/entities"
	"github.com/brooklyncentral/brooklyn-cli/net"
	"github.com/brooklyncentral/brooklyn-cli/terminal"
    "github.com/brooklyncentral/brooklyn-cli/command_metadata"
	"github.com/brooklyncentral/brooklyn-cli/scope"
    "os"
    "fmt"
)

type Entity struct {
	network *net.Network
}

func NewEntity(network *net.Network) (cmd *Entity) {
	cmd = new(Entity)
	cmd.network = network
	return
}

func (cmd *Entity) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "entity",
		Description: "Show the details for an entity",
		Usage:       "BROOKLYN_NAME APP-SCOPE entity",
		Flags:       []cli.Flag{},
	}
}

func (cmd *Entity) Run(scope scope.Scope, c *cli.Context) {
	summary, err := entities.GetEntity(cmd.network, scope.Application, scope.Entity)
    if nil != err {
        fmt.Fprintf(os.Stderr, "Error: %s\n", err)
        os.Exit(1)
    }
	table := terminal.NewTable([]string{"Id", "Name", "Type", "CatalogItemId"})
	table.Add(summary.Id, summary.Name, summary.Type, summary.CatalogItemId)

	table.Print()
}