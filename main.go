package main

import (
	"chart-viewer/cmd"
)

func main() {
	cmd := cmd.NewRootCommand()
	cmd.Execute()
}
