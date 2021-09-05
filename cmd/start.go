package cmd

import (
	log2 "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/DataDrake/waterlog"
	"github.com/DataDrake/waterlog/format"
	"github.com/DataDrake/waterlog/level"
	"github.com/nicolekellydesign/webby-api/server"
)

// Start is the subcommand to start the API service
var Start = cmd.Sub{
	Name:  "start",
	Alias: "s",
	Short: "Start the API service",
	Args:  &StartArgs{},
	Run:   StartFunc,
}

// StartArgs are the args to start the API with
type StartArgs struct {
	Username     string `long:"username" arg:"true" desc:"Database username to use"`
	Password     string `long:"password" arg:"true" desc:"Database password to use"`
	DatabaseName string `long:"database" arg:"true" desc:"Database name to use"`
}

// StartFunc is the function that is run when a user does the Start subcommand
func StartFunc(root *cmd.Root, c *cmd.Sub) {
	// Set up the logger
	logger := waterlog.New(os.Stdout, "webby-cli", log2.Ltime)
	logger.SetLevel(level.Info)
	logger.SetFormat(format.Min)

	// Start our API endpoint listener
	logger.Infoln("Starting the API endpoint listener")
	server := server.Listener{
		Port: 5000,
	}

	go server.Serve()
	logger.Infoln("Now listening on 'localhost:5000'")

	// Wait until told to close
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanup on close
	logger.Println("")
	logger.Goodln("Services shut down successfully")
}
