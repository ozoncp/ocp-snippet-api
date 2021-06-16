package api

import (
	"context"
	"errors"

	"github.com/ozoncp/ocp-snippet-api/internal/models"
	"github.com/ozoncp/ocp-snippet-api/internal/repo"
	desc "github.com/ozoncp/ocp-snippet-api/pkg/ocp-snippet-api"

	// TO BE FIXED: возвращать коды ошибок!
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// const (
// 	errCannotCreateSnippet = "cannot create snippet"
// 	errSnippetNotFound     = "snippet not found"
// )

type api struct {
	desc.UnimplementedOcpSnippetApiServer
	repo repo.Repo
}

func Init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func snippetConvert(snippet *models.Snippet) *desc.Snippet {
	if snippet == nil {
		return nil
	}

	return &desc.Snippet{
		Id:         snippet.Id,
		SolutionId: snippet.SolutionId,
		UserId:     snippet.UserId,
		Text:       snippet.Text,
		Language:   snippet.Language,
	}
}

func (a *api) CreateSnippetV1(ctx context.Context, req *desc.CreateSnippetV1Request) (*desc.CreateSnippetV1Response, error) {
	log.Print("CreateSnippetV1: ", req)

	snippets := []models.Snippet{{
		SolutionId: req.SolutionId,
		UserId:     req.UserId,
		Text:       req.Text,
		Language:   req.Language,
	}}

	if err := a.repo.AddSnippets(ctx, snippets); err != nil {
		return nil, err
	}

	if len(snippets) < 1 {
		return nil, errors.New("Empty snippets received after AddSnippets")
	}

	//err := status.Error(codes.NotFound, errCannotCreateSnippet)
	return &desc.CreateSnippetV1Response{
		Id: snippets[0].Id,
	}, nil
}

func (a *api) DescribeSnippetV1(ctx context.Context, req *desc.DescribeSnippetV1Request) (*desc.DescribeSnippetV1Response, error) {
	log.Print("DescribeSnippetV1: ", req.SnippetId)

	res, err := a.repo.DescribeSnippet(ctx, req.SnippetId)

	if err != nil {
		return nil, err
	}

	return &desc.DescribeSnippetV1Response{
		Snippet: snippetConvert(res),
	}, nil

	// err := status.Error(codes.NotFound, errSnippetNotFound)
	// return nil, err
}

func (a *api) ListSnippetsV1(ctx context.Context, req *desc.ListSnippetsV1Request) (*desc.ListSnippetsV1Response, error) {
	log.Print("ListSnippetsV1: ", req)

	list, err := a.repo.ListSnippets(ctx, req.Limit, req.Offset)

	if err != nil {
		return nil, err
	}

	res := make([]*desc.Snippet, len(list))
	for idx, snippet := range list {
		res[idx] = snippetConvert(&snippet)
	}

	return &desc.ListSnippetsV1Response{
		Snippets: res,
	}, nil

	// err := status.Error(codes.NotFound, errSnippetNotFound)
	// return nil, err
}

func (a *api) RemoveSnippetV1(ctx context.Context, req *desc.RemoveSnippetV1Request) (*desc.RemoveSnippetV1Response, error) {
	log.Print("RemoveSnippetV1: ", req.SnippetId)

	res, err := a.repo.RemoveSnippet(ctx, req.SnippetId)

	return &desc.RemoveSnippetV1Response{
		Removed: res,
	}, err

	// err := status.Error(codes.NotFound, errSnippetNotFound)
	// return nil, err
}

func NewOcpSnippetApi(repo repo.Repo) desc.OcpSnippetApiServer {
	return &api{
		repo: repo,
	}
}
