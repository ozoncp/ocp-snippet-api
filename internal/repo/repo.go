package repo

import (
	"github.com/ozoncp/ocp-snippet-api/internal/snippet" // TO BE REPLACED!
)

type Repo interface {
	// TO BE FIXED: заменить входной параметр на []models.Task, когда ocp-task-api будет доступен для скачивания.
	AddTasks(task []*snippet.Snippet) error
	RemoveTask(taskId uint64) error
	// TO BE FIXED: заменить выходной параметр на *models.Task, когда ocp-task-api будет доступен для скачивания.
	DescribeTask(taskId uint64) (*snippet.Snippet, error)
	// TO BE FIXED: заменить выходной параметр на []models.Task, когда ocp-task-api будет доступен для скачивания.
	ListTasks(limit, offset uint64) ([]*snippet.Snippet, error)
}
