package main

import (
	"database/sql"
	"fmt"
	"os"

	"log"
	"net"
	"net/http"

	"context"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	api "github.com/ozoncp/ocp-snippet-api/internal/api"
	"github.com/ozoncp/ocp-snippet-api/internal/metrics"
	"github.com/ozoncp/ocp-snippet-api/internal/producer"
	"github.com/ozoncp/ocp-snippet-api/internal/repo"
	desc "github.com/ozoncp/ocp-snippet-api/pkg/ocp-snippet-api"

	"github.com/ozoncp/ocp-snippet-api/internal/configuration"
)

// Default values:
const (
	grpcPort = 12345
	httpPort = 54321
)

var (
	grpcEndpoint = fmt.Sprintf("localhost:%d", grpcPort)
	httpEndpoint = fmt.Sprintf("localhost:%d", httpPort)
)

func init() {
	if config := configuration.GetInstance(); config != nil {
		grpcEndpoint = fmt.Sprintf("localhost:%d", config.Grpc.GrpcPort)
		httpEndpoint = fmt.Sprintf("localhost:%d", config.Grpc.HttpPort)
	} else {
		log.Printf("Cannot read config")
	}
}

func getEnv(key string, defaultValue string) string {
	res := os.Getenv(key)
	if len(res) == 0 {
		res = defaultValue
	}
	return res
}

func createDB() *sql.DB {
	dsnPreffix := getEnv("OCP_SNIPPET_API_DB_DSN_PREFFIX", "postgres://")
	host := getEnv("OCP_SNIPPET_API_DB_HOST", "localhost")
	port := getEnv("OCP_SNIPPET_API_DB_PORT", "5432")
	name := getEnv("OCP_SNIPPET_API_DB_NAME", "postgres")
	user := getEnv("OCP_SNIPPET_API_DB_USER", "postgres")
	pswd := getEnv("OCP_SNIPPET_API_DB_PSWD", "")

	dsn := fmt.Sprintf("%s%s:%s@%s", dsnPreffix, user, pswd, host)
	if len(port) > 0 {
		dsn = fmt.Sprintf("%s:%s", dsn, port)
	}
	dsn = fmt.Sprintf("%s/%s", dsn, name)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open DB: %v\n", err)
		return nil
	}

	return db
}

func runMetrics() {

	metrics.RegisterMetrics()
	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(":9100", nil)
	if err != nil {
		panic(err)
	}
}

func run() error {
	listener, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ctx := context.Background()
	s := grpc.NewServer()

	prod, err := producer.NewProducer("ocp-snippet-api")
	if err != nil {
		log.Fatalf("Cannot create prod: %v", err)
	}
	defer prod.Close()

	db := createDB()
	repo := repo.NewRepoDB(db)
	defer db.Close()

	desc.RegisterOcpSnippetApiServer(s, api.NewOcpSnippetApi(repo, prod))

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

	go runMetrics()

	if err := run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("ocp-snippet-api stoped!")
}
