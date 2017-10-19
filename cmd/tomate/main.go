package main

import (
	"flag"
	"os"

	"github.com/mbesancon/tomate"
)

func main() {
	sprintLength := flag.Int("sprint", 30, "Duration of the sprints")
	pauseLength := flag.Int("spause", 5, "Duration of the short pauses")
	longPauseLength := flag.Int("lpause", 10, "Duration of the long pauses")
	nsprints := flag.Int("nsprints", 4, "Number of sprints between long pauses")
	flag.Parse()
	tmt := tomate.New(*sprintLength, *pauseLength, *longPauseLength, *nsprints)
	tmt.Loop(os.Stdout)
}
