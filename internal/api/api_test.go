package api_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"database/sql"
	"database/sql/driver"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/ozoncp/ocp-snippet-api/internal/api"
	"github.com/ozoncp/ocp-snippet-api/internal/models"
	"github.com/ozoncp/ocp-snippet-api/internal/producer"
	"github.com/ozoncp/ocp-snippet-api/internal/repo"
	desc "github.com/ozoncp/ocp-snippet-api/pkg/ocp-snippet-api"
)

var _ = Describe("Api", func() {
	var (
		ctx  context.Context
		db   *sql.DB
		mock sqlmock.Sqlmock

		repoDB    repo.Repo
		apiServer desc.OcpSnippetApiServer
	)

	BeforeEach(func() {
		var err error
		db, mock, err = sqlmock.New()
		if err != nil {
			Fail("Fail to create sqlmock")
		}

		ctx = context.Background()
		repoDB = repo.NewRepoDB(db)

		prod, err := producer.NewProducer("test-ocp-cnippet-api")
		if err != nil {
			Fail("Cannor create producer")
		}

		apiServer = api.NewOcpSnippetApi(repoDB, prod)
	})

	AfterEach(func() {
		db.Close()
	})

	Describe("Create snippet test", func() {
		var (
			req         desc.CreateSnippetV1Request
			expectedRes desc.CreateSnippetV1Response
		)

		It("created successfully", func() {
			req = desc.CreateSnippetV1Request{
				SolutionId: 1,
				Text:       "aa",
				Language:   "ru",
			}
			expectedRes = desc.CreateSnippetV1Response{
				Id: 1,
			}
		})

		It("creating empty snippet", func() {
			req = desc.CreateSnippetV1Request{}
			expectedRes = desc.CreateSnippetV1Response{
				Id: 2,
			}
		})

		AfterEach(func() {
			mock.ExpectQuery("INSERT INTO snippets").
				WithArgs(req.SolutionId, req.Text, req.Language).
				WillReturnRows(
					sqlmock.NewRows([]string{"id"}).AddRow(expectedRes.Id),
				)

			res, err := apiServer.CreateSnippetV1(ctx, &req)

			Expect(err).Should(BeNil())
			Expect(res).ShouldNot(BeNil())
			Expect(res.Id).Should(BeEquivalentTo(expectedRes.Id))
		})
	})

	Describe("Remove snippet test", func() {
		var (
			req         desc.RemoveSnippetV1Request
			expectedRes desc.RemoveSnippetV1Response
			sqlRes      driver.Result
		)

		It("removed successfully", func() {
			req = desc.RemoveSnippetV1Request{
				SnippetId: 2,
			}
			expectedRes = desc.RemoveSnippetV1Response{
				Removed: true,
			}
			sqlRes = sqlmock.NewResult(1, 1)
		})

		It("removing non-existent snippetId", func() {
			req = desc.RemoveSnippetV1Request{
				SnippetId: 0,
			}
			expectedRes = desc.RemoveSnippetV1Response{
				Removed: false,
			}
			sqlRes = sqlmock.NewResult(1, 0)
		})

		AfterEach(func() {
			mock.ExpectExec("DELETE FROM snippets").
				WithArgs(req.SnippetId).
				WillReturnResult(sqlRes)

			res, err := apiServer.RemoveSnippetV1(ctx, &req)

			Expect(err).Should(BeNil())
			Expect(res).ShouldNot(BeNil())
			Expect(res.Removed).Should(BeEquivalentTo(expectedRes.Removed))
		})
	})

	Describe("Describe snippet test", func() {
		var (
			req          desc.DescribeSnippetV1Request
			expectedRes  desc.DescribeSnippetV1Response
			expectedRows *sqlmock.Rows

			columns = []string{"id", "solution_id", "text", "language"}
		)

		It("describe successfully", func() {
			req = desc.DescribeSnippetV1Request{
				SnippetId: 1,
			}
			expectedRes = desc.DescribeSnippetV1Response{
				Snippet: &desc.Snippet{
					Id:         req.SnippetId,
					SolutionId: 1,
					Text:       "aa",
					Language:   "ru",
				},
			}
			expectedRows = sqlmock.NewRows(columns).
				AddRow(expectedRes.Snippet.Id,
					expectedRes.Snippet.SolutionId,
					expectedRes.Snippet.Text,
					expectedRes.Snippet.Language)

			mock.ExpectQuery(`SELECT \* FROM snippets WHERE`).
				WithArgs(req.SnippetId).
				WillReturnRows(expectedRows)

			res, err := apiServer.DescribeSnippetV1(ctx, &req)

			Expect(err).Should(BeNil())
			Expect(res).ShouldNot(BeNil())
			Expect(res.Snippet.Id).Should(BeEquivalentTo(expectedRes.Snippet.Id))
			Expect(res.Snippet.SolutionId).Should(BeEquivalentTo(expectedRes.Snippet.SolutionId))
			Expect(res.Snippet.Text).Should(BeEquivalentTo(expectedRes.Snippet.Text))
			Expect(res.Snippet.Language).Should(BeEquivalentTo(expectedRes.Snippet.Language))
		})

		It("describing non-existent snippetId", func() {
			req = desc.DescribeSnippetV1Request{
				SnippetId: 100,
			}
			expectedRes = desc.DescribeSnippetV1Response{}
			expectedRows = sqlmock.NewRows(columns)

			mock.ExpectQuery(`SELECT \* FROM snippets WHERE`).
				WithArgs(req.SnippetId).
				WillReturnRows(expectedRows)

			res, err := apiServer.DescribeSnippetV1(ctx, &req)

			Expect(err).ShouldNot(BeNil())
			Expect(res).Should(BeNil())
		})
	})

	Describe("List snippet test", func() {
		var (
			req         desc.ListSnippetsV1Request
			expectedRes desc.ListSnippetsV1Response

			rows     *sqlmock.Rows
			snippets []models.Snippet

			columns = []string{"id", "solution_id", "text", "language"}
		)

		BeforeEach(func() {
			snippets = []models.Snippet{
				{Id: 1, SolutionId: 1, Text: "a", Language: "ru"},
				{Id: 2, SolutionId: 2, Text: "b", Language: "en"},
				{Id: 3, SolutionId: 3, Text: "c", Language: "fr"},
				{Id: 4, SolutionId: 4, Text: "d", Language: "ru"},
				{Id: 5, SolutionId: 5, Text: "e", Language: "en"},
				{Id: 6, SolutionId: 6, Text: "f", Language: "ru"},
			}

			rows = sqlmock.NewRows(columns)
		})

		It("correct limit and offet", func() {
			req = desc.ListSnippetsV1Request{
				Limit:  5,
				Offset: 0,
			}
			expectedRes = desc.ListSnippetsV1Response{
				Snippets: []*desc.Snippet{
					{Id: 1, SolutionId: 1, Text: "a", Language: "ru"},
					{Id: 2, SolutionId: 2, Text: "b", Language: "en"},
					{Id: 3, SolutionId: 3, Text: "c", Language: "fr"},
					{Id: 4, SolutionId: 4, Text: "d", Language: "ru"},
					{Id: 5, SolutionId: 5, Text: "e", Language: "en"},
				},
			}

			//for _, snippet := range snippets {
			for i, n := req.Offset, req.Offset+req.Limit; int(i) < len(snippets) && i < n; i++ {
				rows.AddRow(snippets[i].Id,
					snippets[i].SolutionId,
					snippets[i].Text,
					snippets[i].Language)
			}

			mock.ExpectQuery(fmt.Sprintf(`SELECT \* FROM snippets LIMIT %d OFFSET %d`, req.Limit, req.Offset)).
				WillReturnRows(rows)

			res, err := apiServer.ListSnippetsV1(ctx, &req)

			Expect(err).Should(BeNil())
			Expect(res).ShouldNot(BeNil())
			Expect(res.Snippets).Should(Equal(expectedRes.Snippets))
		})

		It("Out of range offset offet", func() {
			req = desc.ListSnippetsV1Request{
				Limit:  15,
				Offset: 100,
			}
			expectedRes = desc.ListSnippetsV1Response{
				Snippets: []*desc.Snippet{},
			}

			mock.ExpectQuery(fmt.Sprintf(`SELECT \* FROM snippets LIMIT %d OFFSET %d`, req.Limit, req.Offset)).
				WillReturnRows(rows)

			res, err := apiServer.ListSnippetsV1(ctx, &req)

			Expect(err).Should(BeNil())
			Expect(res).ShouldNot(BeNil())
			Expect(res.Snippets).Should(Equal(expectedRes.Snippets))
		})
	})
})
