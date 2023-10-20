package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// webPort the port that we listen on for api calls
const webPort = "80"

// Config is the type we'll use as a receiver to share application
// configuration around our app.
type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	// try to connect to RabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	// get service urls
	//app.getServiceURLs()

	// watch service urls
	//go app.watchEtcd()

	log.Println("Starting broker service on port", webPort)

	// define the http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// connect tries to connect to RabbitMQ, and delays between attempts.
// If we can't connect after 5 tries (with increasing delays), return an error
func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbitmq is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			// we have a connection to rabbitmq, so set connection = c and break out of
			// this loop
			connection = c
			break
		}

		if counts > 5 {
			// if we can't connect after five tries, something is wrong...
			fmt.Println(err)
			return nil, err
		}
		fmt.Printf("Backing off for %d seconds...\n", int(math.Pow(float64(counts), 2)))
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}
