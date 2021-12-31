package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"

	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/nicolekellydesign/webby-api/database"
	"github.com/nicolekellydesign/webby-api/server"
)

// AddUserFunc creates a new protected user in the database.
func AddUserFunc(root *cmd.Root, c *cmd.Sub) {
	db, err := database.Connect(dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %s", err)
	}

	args := c.Args.(*AddUserArgs)

	if err = db.AddUser(args.Username, args.Password, true); err != nil {
		log.Fatalf("Error adding user to database: %s\n", err)
	}

	log.Goodln("User added to the database")
}

// ListUsersFunc prints all of the users in the database.
func ListUsersFunc(root *cmd.Root, c *cmd.Sub) {
	db, err := database.Connect(dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %s", err)
	}

	users, err := db.GetUsers()
	if err != nil {
		log.Fatalf("Error getting users from the database: %s\n", err)
	}

	log.Infoln("Users:")
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, user := range users {
		fmt.Fprintf(tw, "\n%s:\n", user.Username)
		fmt.Fprintf(tw, "\tID: %d\n", user.ID)
		fmt.Fprintf(tw, "\tProtected: %t\n", user.Protected)
	}
	tw.Flush()
}

// RemoveUserFunc removes a user from the database.
func RemoveUserFunc(root *cmd.Root, c *cmd.Sub) {
	db, err := database.Connect(dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %s\n", err)
	}

	args := c.Args.(*RemoveUserArgs)

	if err = db.RemoveUser(args.ID); err != nil {
		log.Fatalf("Error removing user from database: %s\n", err)
	}

	log.Goodln("User removed from the database")
}

// ServeFunc opens a database connection and starts serving our HTTP endpoints.
func ServeFunc(root *cmd.Root, c *cmd.Sub) {
	// Start our database connection
	db, err := database.Connect(dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %s\n", err)
	}

	errs := make(chan error, 1)

	// Start our API endpoint listener
	log.Infoln("Starting the API endpoint listener")
	server := server.New(5000, db, log, rootDir, errs)

	go server.Serve()
	log.Infoln("Now listening on 'localhost:5000'")

	// Create a channel to listen for OS signals so we know
	// when to close.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// Wait and block until either an error is received
	// or an OS signal is received telling us to close.
	select {
	case err := <-errs:
		log.Errorf("Error while serving: %s\n", err.Error())
		break
	case <-sc:
		break
	}

	// Cleanup on close
	log.Println("")
	log.Goodln("Services shut down successfully")
}
