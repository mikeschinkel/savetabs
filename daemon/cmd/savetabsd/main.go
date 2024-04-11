//go:build go1.22

// This is an example of implementing the Pet Store from the OpenAPI documentation
// found at:
// https://github.com/OAI/OpenAPI-Specification/blob/master/examples/v3.0/petstore.yaml

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"savetabs/restapi"
	"savetabs/sqlc"
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

var ds sqlc.DataStore

func runServer(port *string) (err error) {
	var swagger *openapi3.T
	var api *restapi.API

	ctx := context.Background()

	ds, err = sqlc.Initialize(ctx, DBFile)
	if err != nil {
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

	err = ui.Initialize(ctx, ds)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error initializing UI: %s", err)
		os.Exit(2)
	}

	restapi.SetErrorTemplate(ui.GetTemplate("error"))
	api = restapi.NewAPI(*port, swagger)
	err = api.ListenAndServe()
end:
	return err
}
