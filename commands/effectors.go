package commands

import (
	"github.com/codegangsta/cli"
	"github.com/brooklyncentral/brooklyn-cli/api/entity_effectors"
	"github.com/brooklyncentral/brooklyn-cli/command_metadata"
	"github.com/brooklyncentral/brooklyn-cli/models"
	"github.com/brooklyncentral/brooklyn-cli/net"
	"github.com/brooklyncentral/brooklyn-cli/terminal"
	"strings"
	"errors"
	"fmt"
	"os"
)

type Effectors struct {
	network *net.Network
}

func NewEffectors(network *net.Network) (cmd *Effectors) {
	cmd = new(Effectors)
	cmd.network = network
	return
}

func (cmd *Effectors) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "effectors",
		Description: "Show the list of effectors for an application and entity",
		Usage:       "BROOKLYN_NAME effectors APPLICATION ENTITY",
		Flags:       []cli.Flag{},
	}
}

func (cmd *Effectors) Run(c *cli.Context) {
	if len(c.Args()) < 3 {
		cmd.listEffectors(c)
	} else {
		err := cmd.invokeEffector(c)
		if nil != err {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}
	}
}

func (cmd *Effectors) listEffectors(c *cli.Context) {
	effectors := entity_effectors.EffectorList(cmd.network, c.Args()[0], c.Args()[1])
	table := terminal.NewTable([]string{"Name", "Description", "Parameters"})
	for _, effector := range effectors {
		var parameters []string
		for _, parameter := range effector.Parameters {
			parameters = append(parameters, parameter.Name)
		}
		table.Add(effector.Name, effector.Description, strings.Join(parameters, ","))
	}
	table.Print()
}

func (cmd *Effectors) invokeEffector(c *cli.Context) error {
	applicationId := c.Args().First()
	entityId := c.Args()[1]
	effectorName := c.Args()[2]
	effectorArgs := c.Args()[3:]

	// get the effector
	effectors := entity_effectors.EffectorList(cmd.network, c.Args()[0], c.Args()[1])
	effector, err := findEffector(effectorName, effectors)
	if nil != err {
		return err
	}

	// get its parameters and check we have an arg for each
	parameters := make([]string, 0)
	for _, parameter := range effector.Parameters {
		parameters = append(parameters, parameter.Name)
	}
	if len(parameters) != len(effectorArgs) {
		return parameterMisMatchError(parameters)
	}

	// invoke it
	_, err = entity_effectors.TriggerEffector(cmd.network, applicationId, entityId, effectorName, parameters, effectorArgs)
	return err
}

func parameterMisMatchError(parameters []string) error {
	pnames := []string{}
	for _, parm := range parameters {
		pnames = append(pnames, parm)
	}
	pnameDesc := strings.Join(pnames, ", ")
	return errors.New(strings.Join([]string{"Parameters not supplied: ", pnameDesc}, ""))
}

func findEffector(name string, effectors []models.EffectorSummary) (models.EffectorSummary, error) {
	for _, effector := range effectors {
		if effector.Name == name {
			return effector, nil
		}
	}
	return models.EffectorSummary{}, errors.New(strings.Join([]string{"Effector not found:", name}, " "))
}