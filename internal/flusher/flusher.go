package flusher

import (
	"github.com/ozoncp/ocp-snippet-api/internal/metrics"
	"github.com/ozoncp/ocp-snippet-api/internal/models"
	"github.com/ozoncp/ocp-snippet-api/internal/repo"
	"github.com/ozoncp/ocp-snippet-api/internal/utils"
)

type Flusher interface {
	Flush(snippets []models.Snippet) []models.Snippet
}

type flusher struct {
	chunkSize uint
	repo      repo.Repo
	publicher metrics.Publisher
}

func (f flusher) Flush(snippets []models.Snippet) []models.Snippet {
	batches, err := utils.SplitSnippetSlice(snippets, f.chunkSize)

	if err != nil {
		f.publicher.PublishFlushing(0)
		return snippets
	}

	res := make([]models.Snippet, 0, len(snippets))

	for _, batch := range batches {
		if err := f.repo.AddSnippets(batch); err != nil {
			res = append(res, batch...)
		}
	}

	f.publicher.PublishFlushing(len(snippets) - len(res))

	if len(res) > 0 {
		return res
	}

	return nil
}

func New(chunkSize uint, repo repo.Repo, publicher metrics.Publisher) Flusher {
	return &flusher{
		chunkSize: chunkSize,
		repo:      repo,
		publicher: publicher,
	}
}

func Init() error {
	return nil
}

func Close() error {
	return nil
}
