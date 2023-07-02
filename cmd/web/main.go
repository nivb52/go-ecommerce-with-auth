package main

import (
	"flag"
	"fmt"
	"go-ecommerce-with-auth/internal/driver"
	"go-ecommerce-with-auth/internal/models"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	// _ "github.com/joho/godotenv/autoload"
)

const version = "1.0.0"
const cssVersion = "1"

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
	templateCache map[string]*template.Template
	version       string
	cssVersion    string
	DB            models.DBModel
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
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		log.Println(":: INFO loading .env file failed:\n", errEnv)
		errEnv = godotenv.Load("./cmd/web/local.env")
		if errEnv != nil {
			log.Println(":: INFO loading local.env file failed:\n", errEnv)
		}
	}

	var cfg config
	// run args
	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application env")
	flag.StringVar(&cfg.api, "api", "http://localhost:4001", "App url")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("DSN"), "Mysql connection string")
	flag.StringVar(&cfg.stripe.key, "stripe_key", os.Getenv("STRIPE_KEY"), "Stripe payments public key")
	flag.Parse()

	// secrets check
	if len(cfg.stripe.secret) < 10 || len(cfg.stripe.key) < 10 {
		log.Fatal("missing Stripe Secret Key")
		os.Exit(2)
	}

	// logs
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	// end logs

	conn, dbErr := driver.OpenDB(cfg.db.dsn)
	if dbErr != nil {
		errorLog.Fatal(":: DB connction failed, error: ", dbErr)
	}

	defer conn.Close()

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
		templateCache: tc,
		version:       version,
		cssVersion:    cssVersion,
		DB:            models.DBModel{DB: conn},
	}

	err := app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal("exiting")
	}
}
