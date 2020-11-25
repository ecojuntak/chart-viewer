package main

import "chart-viewer/cmd/chartviewer"

func main() {
	cmd := chartviewer.NewRootCommand()
	cmd.Execute()
}
