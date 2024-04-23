//go:build go1.22

// This is an example of implementing the Pet Store from the OpenAPI documentation
// found at:
// https://github.com/OAI/OpenAPI-Specification/blob/master/examples/v3.0/petstore.yaml

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/getkin/kin-openapi/openapi3"
	"savetabs/restapi"
	"savetabs/sqlc"
	"savetabs/tasks"
	"savetabs/ui"
)

const (
	defaultPort = "8642"
	DBFile      = "./data/savetabs.db"
)

func main() {
	port := flag.String("port", defaultPort, "Port for test HTTP server")
	flag.Parse()

	err := runServer(port)
	if err != nil {
		log.Fatal(err)
	}
}

func runServer(port *string) (err error) {
	var swagger *openapi3.T
	var api *restapi.API
	var ds sqlc.DataStore
	var stopChan chan os.Signal

	ctx, cancel := context.WithCancel(context.Background())

	ds, err = sqlc.Initialize(ctx, DBFile)
	if err != nil {
		cancel()
		goto end
	}

	swagger, err = restapi.GetSwagger()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	restapi.SetErrorTemplate(ui.GetTemplate("error"))
	api = restapi.NewAPI(restapi.APIArgs{
		Port:    *port,
		Swagger: swagger,
		Views:   ui.NewViews(ds),
	})

	go tasks.BackgroundTask(ctx, tasks.NewCaretaker(ds))

	go func() {
		err := api.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		if err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	// Graceful shutdown
	stopChan = make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan // Wait for interrupt or termination signal

	cancel()
	err = api.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Server shutdown failed: %+v", err)
	}

	slog.Info("API Server and Caretaker stopped gracefully.")
end:
	return err
}
