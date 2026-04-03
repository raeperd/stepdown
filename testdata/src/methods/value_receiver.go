package methods

type Config struct{}

func (c Config) validate() {} // want `function "validate" is called by "load" but declared before it \(stepdown rule\)`

func (c Config) load() {
	c.validate()
}
