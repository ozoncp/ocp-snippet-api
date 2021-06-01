package flusher

import (
	"github.com/ozoncp/ocp-snippet-api/internal/metrics"
	"github.com/ozoncp/ocp-snippet-api/internal/repo"
	models "github.com/ozoncp/ocp-snippet-api/internal/snippet"
	"github.com/ozoncp/ocp-snippet-api/internal/utils"
)

type Flusher interface {
	Flush(task models.Snippets) models.Snippets
}

type flusher struct {
	chunckSize uint
	repo       repo.Repo
	publicher  metrics.Publisher
}

func (f flusher) Flush(task models.Snippets) models.Snippets {
	batches, err := utils.SplitSnippetSlice(&task, f.chunckSize)

	if err != nil {
		f.publicher.PublishFlushing(0)
		return task
	}

	res := make([]*models.Snippet, 0, len(task))

	for _, batch := range batches {
		if err := f.repo.AddSnippets(batch); err != nil {
			res = append(res, batch...)
		}
	}

	f.publicher.PublishFlushing(len(task) - len(res))

	if len(res) > 0 {
		return res
	}

	return nil
}

func New(repo repo.Repo) Flusher {
	return &flusher{
		repo: repo,
	}
}

func Init() error {
	return nil
}

func Close() error {
	return nil
}
