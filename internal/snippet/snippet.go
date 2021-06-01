package models

import (
	"strconv"
)

type Snippet struct {
	Id         uint64
	SolutionId uint64
	UserId     uint64
	Test       string
	Language   string
}
type Snippets = []*Snippet
type SnippetMap = map[uint64]*Snippet

func String(snippetPtr *Snippet) string {
	return "Snippet {" + "\n\t" +
		"Id: " + strconv.Itoa(int(snippetPtr.Id)) + ",\n\t" +
		"SolutionId: " + strconv.Itoa(int(snippetPtr.SolutionId)) + ",\n\t" +
		"UserId: " + strconv.Itoa(int(snippetPtr.UserId)) + ",\n\t" +
		"Test: " + snippetPtr.Test + ",\n\t" +
		"Language: " + snippetPtr.Language + ",\n" +
		"}"
}
