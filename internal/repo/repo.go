package repo

import (
	"context"

	"errors"

	"github.com/ozoncp/ocp-snippet-api/internal/models"

	"database/sql"

	_ "github.com/jackc/pgx/stdlib"

	sq "github.com/Masterminds/squirrel"

	"log"
)

type Repo interface {
	AddSnippets(snippets []models.Snippet) error
	RemoveSnippet(snippetId uint64) (bool, error)
	DescribeSnippet(snippetId uint64) (*models.Snippet, error)
	ListSnippets(limit, offset uint64) ([]models.Snippet, error)
}

const (
	table string = "snippets"
)

type repoDB struct {
	db  *sql.DB
	ctx context.Context
}

func NewRepoDB(db *sql.DB, ctx context.Context) repoDB {
	return repoDB{
		db:  db,
		ctx: ctx,
	}
}

// AddSnippets - добавляет переданные snippet'ы в БД.
// После вставки заполняет поле Id каждого snippet'а.
func (repo repoDB) AddSnippets(snippets []models.Snippet) error {
	if len(snippets) == 0 {
		return errors.New("AddSnippets: Nothing to add")
	}

	// err := repo.db.PingContext(repo.ctx)
	// if err != nil {
	// 	log.Printf("Failed to connect to db: %v\n", err)
	// 	return err
	// }

	query := sq.Insert(table).Columns("solution_id", "user_id", "text", "language").Suffix("RETURNING \"id\"").RunWith(repo.db).PlaceholderFormat(sq.Dollar)
	for _, snippet := range snippets {
		query = query.Values(snippet.SolutionId, snippet.UserId, snippet.Text, snippet.Language)
	}

	rows, err := query.QueryContext(repo.ctx)

	if err != nil {
		log.Printf("Failed to exec query: %v\n", err)
		return err
	}

	for i := 0; i < len(snippets) && rows.Next(); i++ {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			log.Printf("Failed to scan insert query result: %v\n", err)
			return err
		}
		snippets[i].Id = id
	}

	return nil
}

func (repo repoDB) RemoveSnippet(snippetId uint64) (bool, error) {
	query := sq.Delete(table).Where(sq.Eq{"id": snippetId}).RunWith(repo.db).PlaceholderFormat(sq.Dollar)

	res, err := query.ExecContext(repo.ctx)

	if err != nil {
		log.Printf("Failed to exec query: %v\n", err)
		return false, err
	}

	rowsDeleted, _ := res.RowsAffected()
	log.Printf("%d rows deleted!\n", rowsDeleted)

	return rowsDeleted > 0, nil
}

func scanSnippetRow(snippet *models.Snippet, rows *sql.Rows) error {
	return rows.Scan(&snippet.Id, &snippet.SolutionId, &snippet.UserId, &snippet.Text, &snippet.Language)
}

func (repo repoDB) DescribeSnippet(snippetId uint64) (*models.Snippet, error) {
	query := sq.Select("*").From(table).Where(sq.Eq{"id": snippetId}).RunWith(repo.db).PlaceholderFormat(sq.Dollar)

	rows, err := query.QueryContext(repo.ctx)

	if err != nil {
		log.Printf("Failed to exec query: %v\n", err)
		return nil, err
	}

	if !rows.Next() {
		return nil, errors.New("Snippet not found")
	}

	var res models.Snippet
	if err := scanSnippetRow(&res, rows); err != nil {
		log.Printf("Failed to scan rows: %v\n", err)
		return nil, err
	}

	if rows.Next() {
		log.Printf("WARNING: key duplicate in table '%s': %d!\n", table, snippetId)
		//return nil, errors.New("Key duplicate")
	}

	return &res, nil
}

// Если limit == 0 считается, что лимита в запросе нет.
func (repo repoDB) ListSnippets(limit, offset uint64) ([]models.Snippet, error) {
	res := make([]models.Snippet, 0, limit)

	query := sq.Select("*").From(table).Offset(offset).RunWith(repo.db).PlaceholderFormat(sq.Dollar)

	if limit > 0 {
		query = query.Limit(limit)
	}

	rows, err := query.QueryContext(repo.ctx)

	if err != nil {
		log.Printf("Failed to exec query: %v\n", err)
		return nil, err
	}

	for rows.Next() {
		var snippet models.Snippet
		if err := scanSnippetRow(&snippet, rows); err != nil {
			log.Printf("Failed to scan rows: %v\n", err)
			return nil, err
		}

		res = append(res, snippet)
	}

	return res, nil
}
