package storage

import (
	"github.com/ernado/selectel/storage"
	"github.com/earlcherry/selectel-cli/cli"
)

type Containers []storage.ContainerInfo

func (c Containers) Convert() []cli.Tabler {
	var result []cli.Tabler
	for _, container := range c {
		result = append(result, Container(container))
	}
	return result
}

type Container storage.ContainerInfo

func (c Container) Row(n int) []interface{} {
	return []interface{}{n, c.Name, c.Type, c.ObjectCount, c.BytesUsed}
}
func (c Containers) Select(api storage.API) storage.ContainerAPI {
	containers := c.Convert()
	selected:=cli.Table(
		"%v\t%v\t%v\t%v\t%v\t\n",
		[]interface{}{"#", "Container", "Type", "Objects", "Size"},
		[]interface{}{"-", "--------", "----", "--------", "----"},
		containers,
	)
	return api.Container(c[selected].Name)
}
