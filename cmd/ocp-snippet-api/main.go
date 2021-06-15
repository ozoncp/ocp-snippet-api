package main

import (
	"database/sql"
	"fmt"

	"log"
	"net"
	"net/http"

	"context"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	api "github.com/ozoncp/ocp-snippet-api/internal/api"
	"github.com/ozoncp/ocp-snippet-api/internal/repo"
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

func createDB() *sql.DB {
	const (
		dsnPreffix string = "postgres://"
		host       string = "localhost"
		port       int    = 5432
		name       string = "postgres"
		user       string = "postgres"
		pswd       string = "leshiy"
	)

	dsn := fmt.Sprintf("%s%s:%s@%s", dsnPreffix, user, pswd, host)
	if port >= 0 {
		dsn = fmt.Sprintf("%s:%d", dsn, port)
	}
	dsn = fmt.Sprintf("%s/%s", dsn, name)
	fmt.Println(dsn)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open DB: %v\n", err)
		return nil
	}

	return db
}

func run() error {
	listener, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ctx := context.Background()
	s := grpc.NewServer()

	db := createDB()
	repo := repo.NewRepoDB(db, ctx)
	defer db.Close()

	desc.RegisterOcpSnippetApiServer(s, api.NewOcpSnippetApi(repo))

	go func() {
		fmt.Printf("GRPC server listening on %s\n", grpcEndpoint)

		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	runtimeMux := runtime.NewServeMux()
	if err = desc.RegisterOcpSnippetApiHandlerFromEndpoint(ctx, runtimeMux, grpcEndpoint, []grpc.DialOption{grpc.WithInsecure()}); err != nil {
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
