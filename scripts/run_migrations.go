package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

func main() {
	runMigrations()
	runSeed("")
}

func runSeed(seedName string) {
	if seedName == "" {
		fmt.Println("<seed_name> is missing - skipping seed")
		return
	}
	var directory = filepath.Join("./db/seed")

	errEnv := godotenv.Load("./cmd/api/local.env")
	if errEnv != nil {
		log.Fatal("Error loading .env file ", errEnv)
	}

	var dsn string
	flag.StringVar(&dsn, "dsn", os.Getenv("DSN"), "Mysql connection string")
	flag.Parse()
	if len(dsn) < 10 {
		log.Fatal("missing database connection string")
	}

	db, connErr := sql.Open("mysql", dsn)
	if connErr != nil {
		fmt.Println(":: mySQL Driver :: Open error:")
		fmt.Println(connErr)
		return
	}
	defer db.Close()

	filenamePattern := fmt.Sprintf("^%03d-seed.sql$", seedName) // Regex pattern
	var seedFile string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			matched, _ := regexp.MatchString(filenamePattern, info.Name())
			if matched {
				seedFile = path

				if seedFile == "" {
					fmt.Printf("Migration file %03d_.up.sql not found.\n", seedName)
				}

				migrationSQL, readFileErr := os.ReadFile(seedFile)
				if readFileErr != nil {
					fmt.Println("Error reading migration file:", readFileErr)
					// keep to next file
					return nil
				}

				tx, txErr := db.Begin() // Start a transaction
				if txErr != nil {
					fmt.Println("Error starting transaction for migration:", seedName, txErr)
					return txErr
				}

				_, txErr = tx.Exec(string(migrationSQL)) // Execute the SQL within the transaction
				if err != nil {
					fmt.Println("Error executing migration for migration:", seedName, txErr)
					tx.Rollback() // Rollback on error
					return txErr
				}

				txErr = tx.Commit() // Commit if everything is fine
				if txErr != nil {
					fmt.Println("Error committing transaction for migration:", seedName, txErr)
					tx.Rollback()
					return txErr
				}
			}
		}
		return filepath.SkipDir // Stop walking once file is found
	})

	if err != nil {
		fmt.Println("Error walking directory:", err)
		return
	}

	fmt.Printf("Seedes executed successfully %03d.\n", seedName)
}

func runMigrations() {
	var directory = filepath.Join("./db/migrations_sql")

	var migrationNumber = 1
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run run_migration.go <migration_number>")
		fmt.Println("Using default migration number 1, you can abort in 2 seconds")
		time.Sleep(2 * time.Second)
	} else {
		migrationArg, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println("Invalid migration number:", err)
			return
		}
		migrationNumber = migrationArg
	}
	var startMigrationNumber = migrationNumber

	errEnv := godotenv.Load("./cmd/api/local.env")
	if errEnv != nil {
		log.Fatal("Error loading .env file ", errEnv)
	}
	var dsn string
	flag.StringVar(&dsn, "dsn", os.Getenv("DSN"), "Mysql connection string")
	flag.Parse()
	if len(dsn) < 10 {
		log.Fatal("missing database connection string")
	}

	db, connErr := sql.Open("mysql", dsn)
	if connErr != nil {
		fmt.Println(":: mySQL Driver :: Open error:")
		fmt.Println(connErr)
		return
	}
	defer db.Close()

	filenamePattern := fmt.Sprintf("^%03d_.+\\.up\\.sql$", migrationNumber) // Regex pattern
	var migrationFile string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filenamePattern = fmt.Sprintf("^%03d_.+\\.up\\.sql$", migrationNumber) // Regex pattern
			matched, _ := regexp.MatchString(filenamePattern, info.Name())
			if matched {
				migrationFile = path
				//return filepath.SkipDir // Stop walking once file is found

				if migrationFile == "" {
					fmt.Printf("Migration file %03d_.up.sql not found.\n", migrationNumber)
					// keep to next file
					return nil
				}

				migrationSQL, readFileErr := os.ReadFile(migrationFile)
				if readFileErr != nil {
					fmt.Println("Error reading migration file:", readFileErr)
					// keep to next file
					return nil
				}

				tx, txErr := db.Begin() // Start a transaction
				if txErr != nil {
					fmt.Println("Error starting transaction for migration:", migrationNumber, txErr)
					return txErr
				}

				_, txErr = tx.Exec(string(migrationSQL)) // Execute the SQL within the transaction
				if err != nil {
					fmt.Println("Error executing migration for migration:", migrationNumber, txErr)
					tx.Rollback() // Rollback on error
					return txErr
				}

				txErr = tx.Commit() // Commit if everything is fine
				if txErr != nil {
					fmt.Println("Error committing transaction for migration:", migrationNumber, txErr)
					tx.Rollback()
					return txErr
				}
				migrationNumber++
			}
		}
		// keep to next file
		return nil
	})

	if err != nil {
		fmt.Println("Error walking directory:", err)
		return
	}

	fmt.Printf("Migrations executed successfully from %03d to %03d.\n", startMigrationNumber, migrationNumber)
}
