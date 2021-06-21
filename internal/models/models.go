package models

import (
	"strconv"
)

type Snippet struct {
	Id         uint64
	SolutionId uint64
	Text       string
	Language   string
}

func String(snippetPtr *Snippet) string {
	if snippetPtr == nil {
		return ""
	}

	return "Snippet {" + "\n\t" +
		"Id: " + strconv.Itoa(int(snippetPtr.Id)) + ",\n\t" +
		"SolutionId: " + strconv.Itoa(int(snippetPtr.SolutionId)) + ",\n\t" +
		"Text: " + snippetPtr.Text + ",\n\t" +
		"Language: " + snippetPtr.Language + ",\n" +
		"}"
}
