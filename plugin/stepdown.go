// Package plugin registers stepdown as a golangci-lint module plugin.
//
// Import this package for its side effect when building a custom golangci-lint
// binary.
package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"github.com/raeperd/stepdown"
	"golang.org/x/tools/go/analysis"
)

//nolint:gochecknoinits // golangci-lint module plugins register through init.
func init() {
	register.Plugin("stepdown", New)
}

type settings struct {
	Exclusions []string `json:"exclusions"`
}

type stepdownPlugin struct {
	settings settings
}

var _ register.LinterPlugin = (*stepdownPlugin)(nil)

// New creates the golangci-lint module plugin from its configuration.
func New(rawSettings any) (register.LinterPlugin, error) {
	pluginSettings, err := register.DecodeSettings[settings](rawSettings)
	if err != nil {
		return nil, err
	}

	return &stepdownPlugin{settings: pluginSettings}, nil
}

func (p *stepdownPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		stepdown.NewAnalyzer(stepdown.Settings{
			Exclusions: p.settings.Exclusions,
		}),
	}, nil
}

func (*stepdownPlugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
