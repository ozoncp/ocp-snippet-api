package flusher

import (
	"github.com/ozoncp/ocp-snippet-api/internal/metrics"
	"github.com/ozoncp/ocp-snippet-api/internal/repo"
	"github.com/ozoncp/ocp-snippet-api/internal/snippet" // TO BE REPLACED!
	"github.com/ozoncp/ocp-snippet-api/internal/utils"
)

type Flusher interface {
	// TO BE FIXED: заменить входной параметр на []models.Task, когда ocp-task-api будет доступен для скачивания.
	Flush(task []*snippet.Snippet) []*snippet.Snippet
}

type flusher struct {
	chunckSize uint
	repo       repo.Repo
	publicher  metrics.Publisher
}

func (f flusher) Flush(task []*snippet.Snippet) []*snippet.Snippet {
	batches, err := utils.SplitSnippetSlice(&task, f.chunckSize)

	if err != nil {
		return task
	}

	res := make([]*snippet.Snippet, 0, len(task))

	for _, batch := range batches {
		if err := f.repo.AddTasks(batch); err != nil {
			res = append(res, batch...)
		}
	}

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
