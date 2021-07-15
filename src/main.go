package main

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	H "ethapi/handler"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Basic configuration to the logger to output in JSON format
	log.SetFormatter(&log.JSONFormatter{})

	var address, listen string

	// Get the Infura URL from the environment. We support for https and wss
	ia, ok := os.LookupEnv("INFURA_ADDR")
	if !ok {
		log.Fatal("Must pass INFURA_ADDR as environment variable")
	} else {
		address = ia
	}

	// Get the bind address for the webserver from the environment
	la, ok := os.LookupEnv("LISTEN_ADDR")
	if !ok {
		listen = ":9090"
	} else {
		listen = la
	}

	// Get the log level from the environment and configure the logger accordingly
	lvl := os.Getenv("LOG_LEVEL")
	if lvl != "" {
		lvl, err := log.ParseLevel(lvl)
		if err != nil {
			log.Fatal(err)
		}
		log.SetLevel(lvl)
	}

	// Instantiate a ethclient instance to connect to Infura. Fail hard on any error
	client, err := ethclient.Dial(address)
	if err != nil {
		log.Fatal(err)
	}

	// Configure the webserver
	r := mux.NewRouter()
	r.Methods("GET")
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	// Simple two routes setup
	r.Handle("/api/v1/get-block", &H.BlockHandler{
		Client: client,
	})
	r.Handle("/api/v1/get-tx", &H.TxHandler{
		Client: client,
	})

	// Start the webserver or fail hard
	log.Infof("Starting up and binding to %s", listen)
	log.Fatal(http.ListenAndServe(listen, loggedRouter))
}
