package repo

import (
	models "github.com/ozoncp/ocp-snippet-api/internal/snippet"
)

type Repo interface {
	AddSnippets(task models.Snippets) error
	RemoveSnippet(taskId uint64) error
	DescribeSnippet(taskId uint64) (*models.Snippet, error)
	ListSnippets(limit, offset uint64) (models.Snippets, error)
}
