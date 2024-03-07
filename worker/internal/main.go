package internal

import (
	"context"
	"distributed.systems.labs/shared/pkg/communication"
	"distributed.systems.labs/shared/pkg/contracts"
	"distributed.systems.labs/worker/internal/api"
	"distributed.systems.labs/worker/internal/calc"
	"distributed.systems.labs/worker/internal/config"
	"distributed.systems.labs/worker/internal/notify"
	"encoding/json"
	"fmt"

	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func setupCommunicator(comm *communication.RabbitMQCommunicator) error {
	err := comm.DeclareExchange(config.GetRabbitMQTaskExchange())
	if err != nil {
		return err
	}
	err = comm.DeclareQueueAndBind(config.GetRabbitMQTaskQueue(), config.GetRabbitMQTaskExchange())
	if err != nil {
		return err
	}
	err = comm.DeclareExchange(config.GetRabbitMQResultExchange())
	return err
}

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

	connStr, err := config.GetRabbitMQConnStr()
	if err != nil {
		log.Fatalf("failed to get RabbitMQ connection string: %s", err)
	}

	comm, err := communication.InitRabbitMQCommunicator(ctx, connStr, config.GetRabbitMQReconnectTimeout())
	if err != nil {
		log.Fatalf("failed to create communicator: %s", err)
	}
	defer comm.Close()

	err = setupCommunicator(comm)
	if err != nil {
		log.Fatalf("error while setting up communicator: %s", err)
	}

	pubChan := make(chan []byte, 1)
	defer close(pubChan)

	managerNotifier := notify.InitManagerNotifier(ctx, pubChan)
	defer managerNotifier.Close()
	go func() {
		log.Println("starting manager notifier ...")
		managerNotifier.ListenAndNotify()
	}()

	err = comm.RunPublisher(config.GetRabbitMQResultExchange(), pubChan)
	if err != nil {
		log.Fatalf("error while starting publisher: %s", err)
	}
	err = comm.RunConsumer(config.GetRabbitMQTaskQueue(), func(data []byte, logger *log.Logger) error {
		var req contracts.TaskRequest
		err := json.Unmarshal(data, &req)
		if err != nil {
			return fmt.Errorf("error while decoding json in body: %s", err)
		}

		err = req.Validate()
		if err != nil {
			return fmt.Errorf("validation failed with err: %s", err)
		}

		logger.Printf("starting crack for hash: %s, request-id: %s ...", req.ToCrack, req.RequestID)
		go calc.ProcessRequest(ctx, req, managerNotifier.GetResChan())
		return nil
	})

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
