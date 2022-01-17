package main

import (
	"github.com/elliotwms/nosleep/pkg/nosleep"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(nosleep.Analyzer)
}
