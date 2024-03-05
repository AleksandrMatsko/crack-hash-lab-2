package internal

import (
	"context"
	"distributed.systems.labs/worker/internal/api"
	"distributed.systems.labs/worker/internal/config"
	"distributed.systems.labs/worker/internal/notify"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Main() {
	config.ConfigureApp()

	host, port, err := config.GetHostPort()
	if err != nil {
		log.Fatalf("error occured while starting: %s", err)
	}
	log.Printf("configure to listen on http://%s:%s", host, port)

	managerHost, err := config.GetManagerHostAndPort()
	if err != nil {
		log.Fatalf("error occured while starting: %s", err)
	}
	log.Printf("manager %s", managerHost)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	managerNotifier := notify.InitManagerNotifier(ctx)
	defer managerNotifier.Close()
	go func() {
		log.Println("starting manager notifier ...")
		managerNotifier.ListenAndNotify()
	}()

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

	r := api.ConfigureEndpoints(managerNotifier)
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: r,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("listening ...")
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
	log.Println("shutting down")

	ctx, cancelTimeout := context.WithTimeout(ctx, time.Second*10)
	defer cancelTimeout()
	srv.Shutdown(ctx)
}
