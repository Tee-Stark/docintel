package main

import (
	"docintel/pkg/migrate"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	migrateCmd := flag.NewFlagSet("migrate", flag.ExitOnError)
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)

	// Flags for create command
	name := createCmd.String("name", "", "Name of the migration")

	if len(os.Args) < 2 {
		fmt.Println("Expected 'migrate' or 'create' subcommands")
		os.Exit(1)
	}

	config := migrate.NewMigrationConfig()

	switch os.Args[1] {
	case "migrate":
		migrateCmd.Parse(os.Args[2:])
		if err := config.RunMigrations(); err != nil {
			fmt.Printf("Error running migrations: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations completed successfully")

	case "create":
		createCmd.Parse(os.Args[2:])
		if *name == "" {
			fmt.Println("Please provide a name for the migration")
			os.Exit(1)
		}

		// Convert name to snake case and lowercase
		safeName := strings.ToLower(strings.ReplaceAll(*name, " ", "_"))

		if err := config.CreateNewMigration(safeName); err != nil {
			fmt.Printf("Error creating migration: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Println("Expected 'migrate' or 'create' subcommands")
		os.Exit(1)
	}
}
