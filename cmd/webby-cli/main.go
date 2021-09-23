package main

import (
	log2 "log"
	"os"

	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/DataDrake/waterlog"
	"github.com/DataDrake/waterlog/format"
	"github.com/DataDrake/waterlog/level"
)

const (
	envUserKey     = "WEBBY_DB_USER"
	envPasswordKey = "WEBBY_DB_PASSWORD"
	envNameKey     = "WEBBY_DB_NAME"
)

var (
	dbUser     string
	dbPassword string
	dbName     string

	fileLogger *waterlog.WaterLog
	log        *waterlog.WaterLog
	logFile    *os.File
)

// AddUserArgs holds the arguments for the add user command.
type AddUserArgs struct {
	Username string `long:"username" arg:"true" desc:"Username of the user to add"`
	Password string `long:"password" arg:"true" desc:"Password of the user to add"`
}

// RemoveUserArgs holds the arguments for the delete user command.
type RemoveUserArgs struct {
	ID string `long:"id" arg:"true" desc:"ID of the user to remove"`
}

func init() {
	// Set up the loggers
	log = waterlog.New(os.Stdout, "webby-cli", log2.Ltime)
	log.SetLevel(level.Info)
	log.SetFormat(format.Min)

	var err error
	logFile, err = os.Create("webby-api.log")
	if err != nil {
		log.Fatalf("Unable to create or open log file: %s", err)
	}

	fileLogger = waterlog.New(logFile, "webby-db", log2.Ldate|log2.Ltime)
	fileLogger.SetLevel(level.Info)
	fileLogger.SetFormat(format.Un)

	// Get our environment variables
	found := false
	dbUser, found = os.LookupEnv(envUserKey)
	if !found {
		logFile.Close()
		log.Fatalf("required environment variable '%s' not set", envUserKey)
	}

	dbPassword, found = os.LookupEnv(envPasswordKey)
	if !found {
		logFile.Close()
		log.Fatalf("required environment variable '%s' not set", envPasswordKey)
	}

	dbName, found = os.LookupEnv(envNameKey)
	if !found {
		logFile.Close()
		log.Fatalf("required environment variable '%s' not set", envNameKey)
	}
}

func main() {
	root := &cmd.Root{
		Name:  "webby-cli",
		Short: "CLI interface for the Webby API",
	}

	cmd.Register(&cmd.Sub{
		Name:  "adduser",
		Alias: "a",
		Short: "Add a new Webby user to the database",
		Args:  &AddUserArgs{},
		Run:   AddUserFunc,
	})

	cmd.Register(&cmd.Sub{
		Name:  "deluser",
		Alias: "d",
		Short: "Remove a user from the database",
		Args:  &RemoveUserArgs{},
		Run:   RemoveUserFunc,
	})

	cmd.Register(&cmd.Sub{
		Name:  "listusers",
		Alias: "l",
		Short: "List all users in the database",
		Run:   ListUsersFunc,
	})

	cmd.Register(&cmd.Sub{
		Name:  "start",
		Alias: "s",
		Short: "Start the API service",
		Run:   StartFunc,
	})

	defer logFile.Close()

	root.Run()
}
