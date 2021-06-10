package api

import (
	"context"

	desc "github.com/ozoncp/ocp-snippet-api/pkg/ocp-snippet-api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	errCannotCreateSnippet = "cannot create snippet"
	errSnippetNotFound     = "snippet not found"
)

type api struct {
	desc.UnimplementedOcpSnippetApiServer
}

func Init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func (a *api) CreateSnippetV1(ctx context.Context, req *desc.CreateSnippetV1Request) (*desc.CreateSnippetV1Response, error) {
	log.Print("CreateSnippetV1: ", req)

	err := status.Error(codes.NotFound, errCannotCreateSnippet)
	return nil, err
}

func (a *api) DescribeSnippetV1(ctx context.Context, req *desc.DescribeSnippetV1Request) (*desc.DescribeSnippetV1Response, error) {
	log.Print("DescribeSnippetV1: ", req.SnippetId)

	err := status.Error(codes.NotFound, errSnippetNotFound)
	return nil, err
}

func (a *api) ListSnippetsV1(ctx context.Context, req *desc.ListSnippetsV1Request) (*desc.ListSnippetsV1Response, error) {
	log.Print("ListSnippetsV1: ", req)

	err := status.Error(codes.NotFound, errSnippetNotFound)
	return nil, err
}

func (a *api) RemoveSnippetV1(ctx context.Context, req *desc.RemoveSnippetV1Request) (*desc.RemoveSnippetV1Response, error) {
	log.Print("RemoveSnippetV1: ", req.SnippetId)

	err := status.Error(codes.NotFound, errSnippetNotFound)
	return nil, err
}

func NewOcpSnippetApi() desc.OcpSnippetApiServer {
	return &api{}
}
