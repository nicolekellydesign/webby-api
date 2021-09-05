package main

import (
	"github.com/DataDrake/cli-ng/v2/cmd"
	commands "github.com/nicolekellydesign/webby-api/cmd"
)

func main() {
	root := &cmd.Root{
		Name:  "webby-cli",
		Short: "CLI interface for the Webby API",
	}

	cmd.Register(&commands.Start)

	root.Run()
}
