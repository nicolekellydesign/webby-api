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
		log.Fatalf("Unable to connect to the database: %s", err)
	}

	args := c.Args.(*RemoveUserArgs)

	if err = db.RemoveUser(args.ID); err != nil {
		log.Fatalf("Error removing user from database: %s\n", err)
	}

	log.Goodln("User removed from the database")
}

// InitFunc initializes our database schema.
func InitFunc(root *cmd.Root, c *cmd.Sub) {
	db, err := database.Connect(dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %s", err)
	}

	if err = db.InitSchema(); err != nil {
		log.Fatalf("Error initializing tables: %s\n", err)
	}

	log.Goodln("Database schema created successfully")
}

// ServeFunc opens a database connection and starts serving our HTTP endpoints.
func ServeFunc(root *cmd.Root, c *cmd.Sub) {
	// Start our database connection
	db, err := database.Connect(dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %s", err)
	}

	// Start our API endpoint listener
	log.Infoln("Starting the API endpoint listener")
	server := server.New(5000, db, log, uploadDir)

	go server.Serve()
	log.Infoln("Now listening on 'localhost:5000'")

	// Wait until told to close
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanup on close
	log.Println("")
	log.Goodln("Services shut down successfully")
}
