package internal

import (
	"context"
	"distributed.systems.labs/worker/internal/api"
	"distributed.systems.labs/worker/internal/config"
	"distributed.systems.labs/worker/internal/notify"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Main() {
	config.ConfigureApp()

	host := viper.GetString("server.host")
	if host == "" {
		log.Fatalf("no host provided")
	}
	port := viper.GetString("server.port")
	if port == "" {
		log.Fatalf("no port provided")
	}
	log.Printf("configure to listen on http://%s:%s", host, port)
	log.Printf("manager %s:%s", viper.GetString("manager.host"), viper.GetString("manager.port"))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	managerNotifier := notify.InitManagerNotifier(ctx)
	defer managerNotifier.Close()
	go func() {
		log.Println("starting manager notifier ...")
		managerNotifier.ListenAndNotify()
	}()

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
