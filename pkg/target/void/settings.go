package void

import "github.com/sevigo/hokan/pkg/core"

type Configurator struct {
	core.TargetFactory
}

func NewConfigurator() *Configurator {
	return &Configurator{New}
}

func (c *Configurator) Name() string {
	return TargetName
}

func (c *Configurator) Target() core.TargetFactory {
	return c.TargetFactory
}

func (c *Configurator) DefaultConfig() *core.TargetConfig {
	return &core.TargetConfig{
		// always active target for the testing
		Active:      true,
		Name:        TargetName,
		Description: "fake target storage for testing, will print the name of the file",
		Settings: map[string]string{
			"VOID_PREFIX": "",
		},
	}
}

func (c *Configurator) ValidateSettings(settings core.TargetSettings) (bool, error) {
	return true, nil
}
