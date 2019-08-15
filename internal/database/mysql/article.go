package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/golangid/menekel"
	"github.com/sirupsen/logrus"
)

type mysqlArticleRepository struct {
	Conn *sql.DB
}

// NewArticleRepository will create an object that represent the article.Repository interface
func NewArticleRepository(Conn *sql.DB) menekel.ArticleRepository {
	if Conn == nil {
		panic("Database Connections is nil")
	}
	return &mysqlArticleRepository{Conn}
}

func (m *mysqlArticleRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []menekel.Article, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result = make([]menekel.Article, 0)
	for rows.Next() {
		t := menekel.Article{}
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlArticleRepository) Fetch(ctx context.Context, cursor string, num int64) (res []menekel.Article, nextCursor string, err error) {
	qbuilder := squirrel.Select("id", "title", "content", "updated_at", "created_at").From("article")
	qbuilder = qbuilder.OrderBy("id DESC").Limit(uint64(num))

	if cursor != "" {
		decodedCursor, err := strconv.ParseInt(cursor, 10, 64)
		if err != nil && cursor != "" {
			return nil, "", menekel.ErrBadParamInput
		}
		qbuilder = qbuilder.Where(squirrel.Lt{
			"id": decodedCursor,
		})
	}

	query, args, err := qbuilder.ToSql()
	if err != nil {
		return
	}

	res, err = m.fetch(ctx, query, args...)
	if err != nil {
		return nil, "", err
	}

	nextCursor = cursor
	if len(res) > 0 {
		nextCursor = fmt.Sprintf("%d", res[len(res)-1].ID)
	}
	return
}

func (m *mysqlArticleRepository) GetByID(ctx context.Context, id int64) (res menekel.Article, err error) {
	query := `SELECT id,title,content, updated_at, created_at
  						FROM article WHERE ID = ?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return menekel.Article{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, menekel.ErrNotFound
	}

	return
}

func (m *mysqlArticleRepository) GetByTitle(ctx context.Context, title string) (res menekel.Article, err error) {
	query := `SELECT id,title,content, updated_at, created_at
  						FROM article WHERE title = ?`

	list, err := m.fetch(ctx, query, title)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, menekel.ErrNotFound
	}
	return
}

func (m *mysqlArticleRepository) Store(ctx context.Context, a *menekel.Article) (err error) {
	query := `INSERT  article SET title=? , content=? , updated_at=? , created_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, a.Title, a.Content, a.UpdatedAt, a.CreatedAt)
	if err != nil {
		return
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	a.ID = lastID
	return
}

func (m *mysqlArticleRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM article WHERE id = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAfected != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", rowsAfected)
		return
	}

	return
}
func (m *mysqlArticleRepository) Update(ctx context.Context, ar *menekel.Article) (err error) {
	query := `UPDATE article set title=?, content=?, updated_at=? WHERE ID = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, ar.Title, ar.Content, ar.UpdatedAt, ar.ID)
	if err != nil {
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return
	}

	if affect != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", affect)
		return
	}

	return
}
