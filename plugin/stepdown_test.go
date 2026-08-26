package plugin_test

import (
	"path/filepath"
	"testing"

	"github.com/golangci/plugin-module-register/register"
	_ "github.com/raeperd/stepdown/plugin"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestPlugin(t *testing.T) {
	newPlugin, err := register.GetPlugin("stepdown")
	if err != nil {
		t.Fatalf("get registered plugin: %v", err)
	}

	p, err := newPlugin(nil)
	if err != nil {
		t.Fatalf("create plugin: %v", err)
	}

	if got := p.GetLoadMode(); got != register.LoadModeTypesInfo {
		t.Fatalf("load mode = %q, want %q", got, register.LoadModeTypesInfo)
	}

	analyzers, err := p.BuildAnalyzers()
	if err != nil {
		t.Fatalf("build analyzers: %v", err)
	}
	if len(analyzers) != 1 {
		t.Fatalf("analyzer count = %d, want 1", len(analyzers))
	}
	if got := analyzers[0].Name; got != "stepdown" {
		t.Fatalf("analyzer name = %q, want %q", got, "stepdown")
	}
}

func TestPluginSettings(t *testing.T) {
	newPlugin, err := register.GetPlugin("stepdown")
	if err != nil {
		t.Fatalf("get registered plugin: %v", err)
	}

	p, err := newPlugin(map[string]any{
		"exclusions": []string{"init"},
	})
	if err != nil {
		t.Fatalf("create plugin: %v", err)
	}

	analyzers, err := p.BuildAnalyzers()
	if err != nil {
		t.Fatalf("build analyzers: %v", err)
	}

	testdata, err := filepath.Abs(filepath.Join("..", "testdata"))
	if err != nil {
		t.Fatalf("resolve testdata path: %v", err)
	}

	analysistest.Run(t, testdata, analyzers[0], "exclusions")
}

func TestPluginRejectsUnknownSettings(t *testing.T) {
	newPlugin, err := register.GetPlugin("stepdown")
	if err != nil {
		t.Fatalf("get registered plugin: %v", err)
	}

	if _, err := newPlugin(map[string]any{"unknown": true}); err == nil {
		t.Fatal("create plugin with unknown setting: want error")
	}
}
