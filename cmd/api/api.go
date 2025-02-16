package main

import (
	"flag"
	"fmt"
	"go-ecommerce-with-auth/internal/driver"
	"go-ecommerce-with-auth/internal/models"
	"log"
	"net/http"
	"os"
	"time"

	godotenv "github.com/joho/godotenv"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	version  string
	DB       models.DBModel
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

	app.infoLog.Printf("Starting API server in %s mode on port %d", app.config.env, app.config.port)
	return srv.ListenAndServe()
}

func main() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Println(":: INFO loading .env file failed:\n", envErr)
		envErr := godotenv.Load("./cmd/web/local.env")
		if envErr != nil {
			log.Println(":: INFO loading local.env file failed:\n", envErr)
		}
		//  log.Fatal(":: ENV FILE IS MISSING OR WRONG! Exiting")
	}

	var cfg config
	// run args
	flag.IntVar(&cfg.port, "port", 4001, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application env")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("DSN_API"), "Mysql connection string")
	flag.StringVar(&cfg.stripe.key, "stripe_key", os.Getenv("STRIPE_KEY"), "Stripe payments public key")
	flag.StringVar(&cfg.stripe.secret, "stripe_secret", os.Getenv("STRIPE_SECRET"), "Stripe payments secret key")
	flag.Parse()

	// secrets check
	if len(cfg.stripe.secret) < 10 || len(cfg.stripe.key) < 10 {
		log.Fatal("missing Stripe Secret Key")
	}

	// logs
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	conn, connErr := driver.OpenDB(cfg.db.dsn)
	if connErr != nil {
		fmt.Printf("ENV: DSN_API: %s \n", os.Getenv("DSN_API"))
		log.Fatal(":: DB connection Failed! Exiting")
	}
	defer conn.Close()

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  version,
		DB:       models.DBModel{DB: conn},
	}

	err := app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(":: Exiting API")
	}
}
