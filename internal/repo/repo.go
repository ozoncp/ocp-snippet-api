package repo

import (
	"github.com/ozoncp/ocp-snippet-api/internal/models"
)

type Repo interface {
	AddSnippets(task []models.Snippet) error
	RemoveSnippet(taskId uint64) error
	DescribeSnippet(taskId uint64) (*models.Snippet, error)
	ListSnippets(limit, offset uint64) ([]models.Snippet, error)
}
