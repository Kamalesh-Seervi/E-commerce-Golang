package main

import (
	"flag"
	"html/template"
	"log"
	"os"
)

const version = "0.1"
const cssVersion = "1.1"

type config struct {
	port int
	env  string
	api  string
	db   struct {
		dsn string
	}
	stripe struct {
		pubKey    string
		secretKey string
	}
}

type application struct {
	config        config
	infolog       *log.Logger
	errorlog      *log.Logger
	templateCache map[string]*template.Template
	version       string
}

func main() {
	var cfg config
	
	flag.IntVar(&cfg.port, "port", 4000, "server port to listen on")
	flag.StringVar(&cfg.env, "Environment", "Dev", "Type of environment")
	flag.StringVar(&cfg.api, "entrypoint", "http://localhost/4001", "api url")

	flag.Parse()

	cfg.stripe.pubKey = os.Getenv("PUB_key")
	cfg.stripe.secretKey = os.Getenv("Secret_key")

	infoLog := log.New(os.Stdout, "Info\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	tc := make(map[string]*template.Template)

	app := &application{
		config:        cfg,
		infolog:       infoLog,
		errorlog:      errLog,
		templateCache: tc,
		version:       version,
	}

	err := app.server()
	if err != nil {
		panic(err)
	}

}
