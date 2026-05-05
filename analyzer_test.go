package stepdown_test

import (
	"testing"

	"github.com/raeperd/stepdown"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	a := stepdown.NewAnalyzer(stepdown.Settings{})
	analysistest.Run(t, analysistest.TestData(), a, "basic")
	analysistest.Run(t, analysistest.TestData(), a, "valid")
	analysistest.Run(t, analysistest.TestData(), a, "methods")
	analysistest.Run(t, analysistest.TestData(), a, "circular")
	analysistest.Run(t, analysistest.TestData(), a, "calleeorder")
	analysistest.Run(t, analysistest.TestData(), a, "crossstruct")
}

func TestAnalyzerExclusions(t *testing.T) {
	a := stepdown.NewAnalyzer(stepdown.Settings{Exclusions: []string{"init"}})
	analysistest.Run(t, analysistest.TestData(), a, "exclusions")
}

func TestAnalyzerCalleeExclusion(t *testing.T) {
	a := stepdown.NewAnalyzer(stepdown.Settings{Exclusions: []string{"excluded"}})
	analysistest.Run(t, analysistest.TestData(), a, "callee_exclusion")
}

func TestAnalyzerQualifiedExclusion(t *testing.T) {
	a := stepdown.NewAnalyzer(stepdown.Settings{Exclusions: []string{"Server.handle"}})
	analysistest.Run(t, analysistest.TestData(), a, "qualified_exclusion")
}

func TestAnalyzerShortMethodExclusion(t *testing.T) {
	a := stepdown.NewAnalyzer(stepdown.Settings{Exclusions: []string{"handle"}})
	analysistest.Run(t, analysistest.TestData(), a, "short_method_exclusion")
}
