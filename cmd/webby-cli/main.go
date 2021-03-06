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
	envRootKey     = "WEBBY_ROOT"
)

var (
	dbUser     string
	dbPassword string
	dbName     string
	rootDir    string

	log *waterlog.WaterLog
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

	// Get our environment variables
	found := false
	dbUser, found = os.LookupEnv(envUserKey)
	if !found {
		log.Fatalf("required environment variable '%s' not set\n", envUserKey)
	}

	dbPassword, found = os.LookupEnv(envPasswordKey)
	if !found {
		log.Fatalf("required environment variable '%s' not set\n", envPasswordKey)
	}

	dbName, found = os.LookupEnv(envNameKey)
	if !found {
		log.Fatalf("required environment variable '%s' not set\n", envNameKey)
	}

	rootDir, found = os.LookupEnv(envRootKey)
	if !found {
		log.Fatalf("required environment variable '%s' not set\n", envRootKey)
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
		Name:  "serve",
		Alias: "s",
		Short: "Start serving the API endpoints",
		Run:   ServeFunc,
	})

	root.Run()
}
