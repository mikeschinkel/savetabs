package restapi

//go:generate oapi-codegen -package restapi -generate "std-http"	-o api_server.go 	./openapi.yaml
//go:generate oapi-codegen -package restapi -generate "types"	 		-o api_types.go 	./openapi.yaml
//go:generate oapi-codegen -package restapi -generate "spec"	 		-o api_spec.go 		./openapi.yaml
//go:generate oapi-codegen -package restapi -generate "client"	  -o api_client.go 	./openapi.yaml

