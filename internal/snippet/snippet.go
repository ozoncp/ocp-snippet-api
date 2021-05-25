package snippet

import (
	"strconv"
)

type Snippet struct {
	UserId   uint64
	Test     string
	Language string
}
type SnippetSlice = []*Snippet
type SnippetMap = map[uint64]*Snippet

func String(snippetPtr *Snippet) string {
	return "Snippet{\n\tUserId: " + strconv.Itoa(int(snippetPtr.UserId)) + ",\n\tTest: " + snippetPtr.Test + ",\n\tLanguage: " + snippetPtr.Language + "\n}"
}

func CompareSnippets(l *Snippet, r *Snippet) bool {
	return l == r || (l.UserId == r.UserId &&
		l.Test == r.Test &&
		l.Language == r.Language)
}
