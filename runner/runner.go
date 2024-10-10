package runner

import (
	"container-log-backend/config"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func Run(router http.Handler, serverConfiguration *config.ServerConfiguration) (err error) {
	shutdown := make(chan error)
	go doShutdownOnSignal(shutdown)

	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(serverConfiguration.Port),
		Handler: router,
	}
	go func(httpServer *http.Server, configuration *config.ServerConfiguration) {
		if configuration.SSL == 0 {
			// HTTP Server
			err = httpServer.ListenAndServe()
		} else {
			// HTTPS Server
			err = httpServer.ListenAndServeTLS(configuration.CertPath, configuration.KeyPath)
		}
		doShutdown(shutdown, err)
	}(httpServer, serverConfiguration)

	err = <-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return httpServer.Shutdown(ctx)
}

func doShutdownOnSignal(shutdown chan<- error) {
	onSignal := make(chan os.Signal, 1)
	signal.Notify(onSignal, os.Interrupt, syscall.SIGTERM)
	sig := <-onSignal
	doShutdown(shutdown, fmt.Errorf("received signal %s", sig))
}

func doShutdown(shutdown chan<- error, err error) {
	select {
	case shutdown <- err:
	default:
		// If there is no one listening on the shutdown channel, then the
		// shutdown is already initiated and we can ignore these errors.
	}
}
