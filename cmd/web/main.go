package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"go-ecommerce-with-auth/internal/driver"
	"go-ecommerce-with-auth/internal/models"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
	// _ "github.com/joho/godotenv/autoload"
)

const version = "1.0.0"
const cssVersion = "1"

var sessionManager *scs.SessionManager

type config struct {
	port int
	env  string
	api  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
}

type application struct {
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	debugLog      *log.Logger
	templateCache map[string]*template.Template
	version       string
	cssVersion    string
	DB            models.DBModel
	Session       *scs.SessionManager
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Printf("Starting HTTP server in %s mode on port %d", app.config.env, app.config.port)
	return srv.ListenAndServe()
}

func main() {
	// DEBUG ENV WITH: printCurrentFolderContent()
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	errEnv := godotenv.Load("./.env")
	if errEnv != nil {
		log.Println(":: PRE RUNNING: loading .env file failed:\n", errEnv)
		errEnv = godotenv.Load("./cmd/web/.env")
		if errEnv != nil {
			log.Println(":: PRE RUNNING: loading local.env file failed:\n", errEnv)
		}
	}
	var cfg config
	// run args
	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", env, "Application Environment variable")
	flag.StringVar(&cfg.api, "api", "http://localhost:4001", "App url")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("DSN"), "Mysql connection string")
	flag.StringVar(&cfg.stripe.key, "stripe_key", os.Getenv("STRIPE_KEY"), "Stripe payments public key")
	flag.StringVar(&cfg.stripe.secret, "stripe_secret", os.Getenv("STRIPE_SECRET"), "Stripe payments secret key")
	flag.Parse()

	// secrets check
	if len(cfg.stripe.secret) < 5 || len(cfg.stripe.key) < 5 {
		log.Fatal("missing Stripe Secret or Public Key, exiting ... ")
		os.Exit(2)
	}

	gob.Register(map[string]any{})

	// logs
	infoLog := log.New(os.Stdout, "::INFO :\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "::ERROR :\t", log.Ldate|log.Ltime|log.Lshortfile)
	var debugLog *log.Logger
	if cfg.env == "development" || cfg.env == "dev" {
		debugLog = log.New(os.Stdout, "DEBUG: ", log.LstdFlags)
	} else {
		// Create a logger with a discard output (i.e., no logging)
		debugLog = log.New(ioutil.Discard, "", log.LstdFlags)
	}
	// end logs

	conn, dbErr := driver.OpenDB(cfg.db.dsn)
	if dbErr != nil {
		errorLog.Fatal(":: DB connction failed, error: ", dbErr)
	}

	defer conn.Close()

	//set up session
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	conn, connErr := driver.OpenDB(cfg.db.dsn)
	if connErr != nil {
		log.Fatal(":: DB connecition Failed! Exiting")
	}
	defer conn.Close()
	tc := make(map[string]*template.Template)

	app := &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		debugLog:      debugLog,
		templateCache: tc,
		version:       version,
		cssVersion:    cssVersion,
		DB:            models.DBModel{DB: conn},
		Session:       sessionManager,
	}

	err := app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal("exiting")
	}
}

func printCurrentFolderContent() {
	// Get the current directory
	currentDir := "./"

	// Read the directory content
	files, err := ioutil.ReadDir(currentDir)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate over the files and directories
	for _, file := range files {
		fmt.Println(file.Name())
	}

	filePath := "./.env"

	// Read the file contents
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the file contents to a string and print it
	fileContent := string(content)
	fmt.Println(fileContent)
}
