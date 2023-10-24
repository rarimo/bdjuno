package v3

import (
	v4 "github.com/forbole/juno/v4/cmd/migrate/v4"

	"github.com/rarimo/bdjuno/modules/actions"
)

type Config struct {
	v4.Config `yaml:"-,inline"`

	// The following are there to support modules which config are present if they are enabled

	Actions *actions.Config `yaml:"actions"`
}
