package internal

import (
	"context"
	"distributed.systems.labs/shared/pkg/communication"
	"distributed.systems.labs/shared/pkg/contracts"
	"distributed.systems.labs/worker/internal/calc"
	"distributed.systems.labs/worker/internal/config"
	"distributed.systems.labs/worker/internal/notify"
	"encoding/json"
	"fmt"

	"log"
	"os"
	"os/signal"
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
	log.Println("shutting down")
}
