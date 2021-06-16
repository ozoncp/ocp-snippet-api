package repo

import (
	"context"
	"os"

	"errors"

	"github.com/ozoncp/ocp-snippet-api/internal/models"
	"github.com/ozoncp/ocp-snippet-api/internal/utils"
	"github.com/rs/zerolog"

	"database/sql"

	_ "github.com/jackc/pgx/stdlib"

	sq "github.com/Masterminds/squirrel"
)

type Repo interface {
	AddSnippets(ctx context.Context, snippets []models.Snippet) error
	RemoveSnippet(ctx context.Context, snippetId uint64) (bool, error)
	DescribeSnippet(ctx context.Context, snippetId uint64) (*models.Snippet, error)
	ListSnippets(ctx context.Context, limit, offset uint64) ([]models.Snippet, error)
}

const (
	table     string = "snippets"
	chunkSize uint   = 10 // TO BE FIXED: надо как-то параметризировать
)

var (
	log = zerolog.New(os.Stdout)
)

type repoDB struct {
	db *sql.DB
}

func NewRepoDB(db *sql.DB) repoDB {
	return repoDB{
		db: db,
	}
}

func pingContext(ctx context.Context, repo repoDB) error {
	err := repo.db.PingContext(ctx)
	if err != nil {
		log.Warn().Msgf("Failed to connect to db: %v\n", err)
	}
	return err
}

// AddSnippets - добавляет переданные snippet'ы в БД.
// После вставки заполняет поле Id каждого snippet'а.
func (repo repoDB) AddSnippets(ctx context.Context, snippets []models.Snippet) error {
	if len(snippets) == 0 {
		return errors.New("AddSnippets: Nothing to add")
	}

	if err := pingContext(ctx, repo); err != nil {
		return err
	}

	batches, err := utils.SplitSnippetSlice(snippets, chunkSize)

	if err != nil {
		log.Warn().Msgf("Failed to split snippets: %v\n", err)
		return err
	}

	for _, batch := range batches {
		query := sq.Insert(table).Columns("solution_id", "user_id", "text", "language").Suffix("RETURNING \"id\"").RunWith(repo.db).PlaceholderFormat(sq.Dollar)
		for _, snippet := range batch {
			query = query.Values(snippet.SolutionId, snippet.UserId, snippet.Text, snippet.Language)
		}

		rows, err := query.QueryContext(ctx)

		if err != nil {
			log.Warn().Msgf("Failed to exec query: %v\n", err)
			return err
		}

		for i := 0; i < len(batch) && rows.Next(); i++ {
			var id uint64
			if err := rows.Scan(&id); err != nil {
				log.Warn().Msgf("Failed to scan insert query result: %v\n", err)
				return err
			}
			batch[i].Id = id
		}
	}

	return nil
}

func (repo repoDB) RemoveSnippet(ctx context.Context, snippetId uint64) (bool, error) {
	if err := pingContext(ctx, repo); err != nil {
		return false, err
	}

	query := sq.Delete(table).Where(sq.Eq{"id": snippetId}).RunWith(repo.db).PlaceholderFormat(sq.Dollar)

	res, err := query.ExecContext(ctx)

	if err != nil {
		log.Warn().Msgf("Failed to exec query: %v\n", err)
		return false, err
	}

	rowsDeleted, _ := res.RowsAffected()
	log.Info().Msgf("%d rows deleted!\n", rowsDeleted)

	return rowsDeleted > 0, nil
}

func scanSnippetRow(snippet *models.Snippet, rows *sql.Rows) error {
	return rows.Scan(&snippet.Id, &snippet.SolutionId, &snippet.UserId, &snippet.Text, &snippet.Language)
}

func (repo repoDB) DescribeSnippet(ctx context.Context, snippetId uint64) (*models.Snippet, error) {
	if err := pingContext(ctx, repo); err != nil {
		return nil, err
	}

	query := sq.Select("*").From(table).Where(sq.Eq{"id": snippetId}).RunWith(repo.db).PlaceholderFormat(sq.Dollar)

	rows, err := query.QueryContext(ctx)

	if err != nil {
		log.Warn().Msgf("Failed to exec query: %v\n", err)
		return nil, err
	}

	if !rows.Next() {
		return nil, errors.New("snippet not found")
	}

	var res models.Snippet
	if err := scanSnippetRow(&res, rows); err != nil {
		log.Warn().Msgf("Failed to scan rows: %v\n", err)
		return nil, err
	}

	if rows.Next() {
		log.Warn().Msgf("key duplicate in table '%s': %d!\n", table, snippetId)
		//return nil, errors.New("Key duplicate")
	}

	return &res, nil
}

// Если limit == 0 считается, что лимита в запросе нет.
func (repo repoDB) ListSnippets(ctx context.Context, limit, offset uint64) ([]models.Snippet, error) {
	if err := pingContext(ctx, repo); err != nil {
		return nil, err
	}

	res := make([]models.Snippet, 0, limit)

	query := sq.Select("*").From(table).Offset(offset).RunWith(repo.db).PlaceholderFormat(sq.Dollar)

	if limit > 0 {
		query = query.Limit(limit)
	}

	rows, err := query.QueryContext(ctx)

	if err != nil {
		log.Warn().Msgf("Failed to exec query: %v\n", err)
		return nil, err
	}

	for rows.Next() {
		var snippet models.Snippet
		if err := scanSnippetRow(&snippet, rows); err != nil {
			log.Warn().Msgf("Failed to scan rows: %v\n", err)
			return nil, err
		}

		res = append(res, snippet)
	}

	return res, nil
}
