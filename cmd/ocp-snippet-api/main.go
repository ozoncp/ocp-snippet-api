package main

import (
	"fmt"

	"log"
	"net"
	"net/http"

	"context"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	api "github.com/ozoncp/ocp-snippet-api/internal/api"
	desc "github.com/ozoncp/ocp-snippet-api/pkg/ocp-snippet-api"
)

const (
	grpcPort = 12345
	httpPort = 54321
)

var (
	grpcEndpoint = fmt.Sprintf("localhost:%d", grpcPort)
	httpEndpoint = fmt.Sprintf("localhost:%d", httpPort)
)

func run() error {
	listen, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	desc.RegisterOcpSnippetApiServer(s, api.NewOcpSnippetApi())

	go func() {
		fmt.Printf("GRPC server listening on %s\n", grpcEndpoint)

		if err := s.Serve(listen); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	runtimeMux := runtime.NewServeMux()
	if err = desc.RegisterOcpSnippetApiHandlerFromEndpoint(context.Background(), runtimeMux, grpcEndpoint, []grpc.DialOption{grpc.WithInsecure()}); err != nil {
		log.Fatalf("failed to register handler: %v", err)
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/", runtimeMux)

	fmt.Printf("HTTP server listening on %s\n", httpEndpoint)
	return http.ListenAndServe(httpEndpoint, httpMux)
}

func main() {
	fmt.Println("ocp-snippet-api by Oleg Usov")

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
