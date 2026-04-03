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
}
