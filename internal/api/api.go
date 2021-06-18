package api

import (
	"context"
	"errors"

	"github.com/ozoncp/ocp-snippet-api/internal/metrics"
	"github.com/ozoncp/ocp-snippet-api/internal/models"
	"github.com/ozoncp/ocp-snippet-api/internal/producer"
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
	prod producer.Producer
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
		return nil, errors.New("empty snippets received after AddSnippets")
	}

	res := &desc.CreateSnippetV1Response{
		Id: snippets[0].Id,
	}

	a.prod.SendMessage(producer.SnippetEvent{
		Type: producer.Created,
		Body: map[string]interface{}{
			"Id": res.Id,
		},
	})
	metrics.IncrementSuccessfulCreate(len(snippets))

	//err := status.Error(codes.NotFound, errCannotCreateSnippet)
	return res, nil
}

func (a *api) MultiCreateSnippetV1(ctx context.Context, req *desc.MultiCreateSnippetV1Request) (*desc.MultiCreateSnippetV1Response, error) {
	log.Print("CreateSnippetV1: ", req)

	snippets := make([]models.Snippet, len(req.Snippets))
	for idx, snippet := range req.Snippets { // Корявенько выглядит...
		snippets[idx].SolutionId = snippet.SolutionId
		snippets[idx].UserId = snippet.UserId
		snippets[idx].Text = snippet.Text
		snippets[idx].Language = snippet.Language
	}

	if err := a.repo.AddSnippets(ctx, snippets); err != nil {
		return nil, err
	}

	if len(snippets) < 1 {
		return nil, errors.New("empty snippets received after AddSnippets")
	}

	metrics.IncrementSuccessfulCreate(len(snippets))

	res := make([]uint64, len(snippets))
	for idx, snippet := range snippets {
		res[idx] = snippet.Id
	}

	a.prod.SendMessage(producer.SnippetEvent{
		Type: producer.Created,
		Body: map[string]interface{}{
			"Ids": res,
		},
	})

	//err := status.Error(codes.NotFound, errCannotCreateSnippet)
	return &desc.MultiCreateSnippetV1Response{
		Ids: res,
	}, nil
}

func (a *api) DescribeSnippetV1(ctx context.Context, req *desc.DescribeSnippetV1Request) (*desc.DescribeSnippetV1Response, error) {
	log.Print("DescribeSnippetV1: ", req.SnippetId)

	res, err := a.repo.DescribeSnippet(ctx, req.SnippetId)

	if err != nil {
		return nil, err
	}

	metrics.IncrementSuccessfulRead(1)

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

	metrics.IncrementSuccessfulRead(len(list))

	return &desc.ListSnippetsV1Response{
		Snippets: res,
	}, nil

	// err := status.Error(codes.NotFound, errSnippetNotFound)
	// return nil, err
}

func (a *api) RemoveSnippetV1(ctx context.Context, req *desc.RemoveSnippetV1Request) (*desc.RemoveSnippetV1Response, error) {
	log.Print("RemoveSnippetV1: ", req.SnippetId)

	res, err := a.repo.RemoveSnippet(ctx, req.SnippetId)

	a.prod.SendMessage(producer.SnippetEvent{
		Type: producer.Created,
		Body: map[string]interface{}{
			"Id":      req.SnippetId,
			"Removed": res,
		},
	})

	var deletedCnt int
	if res {
		deletedCnt = 1
	}
	metrics.IncrementSuccessfulDelete(deletedCnt)

	return &desc.RemoveSnippetV1Response{
		Removed: res,
	}, err

	// err := status.Error(codes.NotFound, errSnippetNotFound)
	// return nil, err
}

func (a *api) UpdateSnippetV1(ctx context.Context, req *desc.UpdateSnippetV1Request) (*desc.UpdateSnippetV1Response, error) {
	log.Print("UpdateSnippetV1: ", req)

	if req == nil {
		return nil, errors.New("empty snippet received fo update")
	}

	res, err := a.repo.UpdateSnippet(ctx, models.Snippet{
		Id:         req.Id,
		SolutionId: req.SolutionId,
		UserId:     req.UserId,
		Text:       req.Text,
		Language:   req.Language,
	})

	var updatedCnt int
	if res {
		updatedCnt = 1
	}
	metrics.IncrementSuccessfulUpdate(updatedCnt)

	return &desc.UpdateSnippetV1Response{
		Updated: res,
	}, err
}

func NewOcpSnippetApi(repo repo.Repo, prod producer.Producer) desc.OcpSnippetApiServer {
	return &api{
		repo: repo,
		prod: prod,
	}
}
