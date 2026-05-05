package methods

type Config struct{}

func (c Config) validate() {} // want `function "Config.validate" is called by "Config.load" but declared before it \(stepdown rule\)`

func (c Config) load() {
	c.validate()
}
