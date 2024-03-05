package internal

import (
	"context"
	"distributed.systems.labs/manager/internal/api"
	"distributed.systems.labs/manager/internal/config"
	"distributed.systems.labs/manager/internal/storage"
	"distributed.systems.labs/shared/pkg/alphabet"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func prepareAlphabetRunes() []rune {
	runes := make([]rune, 0)
	for r := 'a'; r <= 'z'; r++ {
		runes = append(runes, r)
	}
	for r := '0'; r <= '9'; r++ {
		runes = append(runes, r)
	}
	return runes
}

func Main() {
	config.ConfigureApp()
	host, port, err := config.GetHostPort()
	if err != nil {
		log.Fatalf("error occured while starting: %s", err)
	}
	log.Printf("configure to listen on http://%s:%s", host, port)
	log.Printf("workers %s", viper.GetString("workers.list"))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var store storage.Storage
	connStr, err := config.GetMongoConnStr()
	if err != nil {
		log.Printf("failed to get mongo connection string: %s", err)
		log.Printf("using in memory storage...")
		store = storage.InitInMemoryStorage(ctx)
	} else {
		store, err = storage.InitMongoStorage(ctx, connStr)
		if err != nil {
			log.Fatalf("failed to connect to mongo: %s", err)
		}
		defer store.Close()
		log.Println("successfully connected to mongodb")
	}

	mqConnStr, err := config.GetRabbitMQConnStr()
	if err != nil {
		log.Fatalf("failed to get RabbitMQ connection string: %s", err)
	}
	connection, err := amqp.Dial(mqConnStr)
	if err != nil {
		log.Fatalf("failed to establish connection with RabbitMQ: %s", err)
	}
	defer connection.Close()
	log.Println("successfully connected to RabbitMQ")

	A := alphabet.InitAlphabet(prepareAlphabetRunes())
	log.Printf("alphabet: '%s'", A.ToOneLine())

	r := api.ConfigureEndpoints(store, A)
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: r,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
	log.Println("shutting down")

	ctx, cancelByTimeout := context.WithTimeout(ctx, time.Second*10)
	defer cancelByTimeout()
	srv.Shutdown(ctx)
}
