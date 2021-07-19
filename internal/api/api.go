package api

import (
	"context"
	"fmt"

	"github.com/ozoncp/ocp-snippet-api/internal/metrics"
	"github.com/ozoncp/ocp-snippet-api/internal/models"
	"github.com/ozoncp/ocp-snippet-api/internal/producer"
	"github.com/ozoncp/ocp-snippet-api/internal/repo"
	desc "github.com/ozoncp/ocp-snippet-api/pkg/ocp-snippet-api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	errNilRequest            = "nil request received"
	errConverter             = "converter error"
	errCannotDescribeSnippet = "cannot describe snippet"
	errEmptySnippetsAfterAdd = "empty snippets received after AddSnippets"
	errCannotListSnippet     = "cannot list snippets"
	errCannotRemoveSnippet   = "cannot remove snippet"
	errCannotUpdateSnippet   = "cannot update snippet"
	errCannotRestoreSnippet  = "cannot restore snippet"
)

type api struct {
	desc.UnimplementedOcpSnippetApiServer
	repo repo.Repo
	prod producer.Producer
}

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func snippetConvert(snippet *models.Snippet) *desc.Snippet {
	if snippet == nil {
		return nil
	}

	return &desc.Snippet{
		Id:         snippet.Id,
		SolutionId: snippet.SolutionId,
		Text:       snippet.Text,
		Language:   snippet.Language,
	}
}
func createSnippetRequestConverter(snippet *desc.CreateSnippetV1Request) *models.Snippet {
	if snippet == nil {
		return nil
	}

	return &models.Snippet{
		SolutionId: snippet.SolutionId,
		Text:       snippet.Text,
		Language:   snippet.Language,
	}
}

func (a *api) CreateSnippetV1(ctx context.Context, req *desc.CreateSnippetV1Request) (*desc.CreateSnippetV1Response, error) {
	log.Print("CreateSnippetV1: ", req)

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, errNilRequest)
	}

	snippet := createSnippetRequestConverter(req)
	if snippet == nil {
		return nil, status.Error(codes.FailedPrecondition, errConverter)
	}

	snippets := []models.Snippet{*snippet}

	if err := a.repo.AddSnippets(ctx, snippets); err != nil {
		return nil, status.Error(codes.DataLoss, err.Error())
	}

	if len(snippets) < 1 {
		return nil, status.Error(codes.DataLoss, errEmptySnippetsAfterAdd)
	}

	fmt.Println(snippet)
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

	return res, nil
}

func (a *api) MultiCreateSnippetV1(ctx context.Context, req *desc.MultiCreateSnippetV1Request) (*desc.MultiCreateSnippetV1Response, error) {
	log.Print("CreateSnippetV1: ", req)

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, errNilRequest)
	}

	snippets := make([]models.Snippet, len(req.Snippets))
	for idx, reqSnippet := range req.Snippets {
		snippet := createSnippetRequestConverter(reqSnippet)
		if snippet == nil {
			log.Warn().Msgf("cannot convert request to models.Snippet (wiil be skipped): %v", reqSnippet)
		} else {
			snippets[idx] = *snippet
		}
	}

	if err := a.repo.AddSnippets(ctx, snippets); err != nil {
		return nil, err
	}

	if len(snippets) < 1 {
		return nil, status.Error(codes.DataLoss, errEmptySnippetsAfterAdd)
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

	return &desc.MultiCreateSnippetV1Response{
		Ids: res,
	}, nil
}

func (a *api) DescribeSnippetV1(ctx context.Context, req *desc.DescribeSnippetV1Request) (*desc.DescribeSnippetV1Response, error) {
	log.Print("DescribeSnippetV1: ", req.SnippetId)

	res, err := a.repo.DescribeSnippet(ctx, req.SnippetId)

	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("%s: %s", errCannotDescribeSnippet, err.Error())) //err
	}

	metrics.IncrementSuccessfulRead(1)

	return &desc.DescribeSnippetV1Response{
		Snippet: snippetConvert(res),
	}, nil
}

func (a *api) ListSnippetsV1(ctx context.Context, req *desc.ListSnippetsV1Request) (*desc.ListSnippetsV1Response, error) {
	log.Print("ListSnippetsV1: ", req)

	list, err := a.repo.ListSnippets(ctx, req.Limit, req.Offset)

	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("%s: %s", errCannotListSnippet, err.Error())) //err
	}

	res := make([]*desc.Snippet, len(list))
	for idx, snippet := range list {
		res[idx] = snippetConvert(&snippet)
	}

	metrics.IncrementSuccessfulRead(len(list))

	return &desc.ListSnippetsV1Response{
		Snippets: res,
	}, nil
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

	if err != nil {
		err = status.Error(codes.Aborted, fmt.Sprintf("%s: %s", errCannotRemoveSnippet, err.Error()))
	}

	return &desc.RemoveSnippetV1Response{
		Removed: res,
	}, err
}

func (a *api) RestoreSnippetV1(ctx context.Context, req *desc.RestoreSnippetV1Request) (*desc.RestoreSnippetV1Response, error) {
	log.Print("RestoreSnippetV1: ", req.SnippetId)

	res, err := a.repo.RestoreSnippet(ctx, req.SnippetId)

	a.prod.SendMessage(producer.SnippetEvent{
		Type: producer.Created,
		Body: map[string]interface{}{
			"Id":       req.SnippetId,
			"Restored": res,
		},
	})

	var deletedCnt int
	if res {
		deletedCnt = 1
	}
	metrics.IncrementSuccessfulDelete(deletedCnt)

	if err != nil {
		err = status.Error(codes.Aborted, fmt.Sprintf("%s: %s", errCannotRestoreSnippet, err.Error()))
	}

	return &desc.RestoreSnippetV1Response{
		Restored: res,
	}, err
}

func (a *api) UpdateSnippetV1(ctx context.Context, req *desc.UpdateSnippetV1Request) (*desc.UpdateSnippetV1Response, error) {
	log.Print("UpdateSnippetV1: ", req)

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, errNilRequest)
	}

	res, err := a.repo.UpdateSnippet(ctx, models.Snippet{
		Id:         req.Id,
		SolutionId: req.SolutionId,
		Text:       req.Text,
		Language:   req.Language,
	})

	var updatedCnt int
	if res {
		updatedCnt = 1
	}
	metrics.IncrementSuccessfulUpdate(updatedCnt)

	if err != nil {
		err = status.Error(codes.Aborted, fmt.Sprintf("%s: %s", errCannotUpdateSnippet, err.Error()))
	}

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
