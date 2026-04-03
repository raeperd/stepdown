package main

import (
	"github.com/raeperd/stepdown"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(stepdown.NewAnalyzer(stepdown.Settings{}))
}
