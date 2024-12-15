package service

import (
	"context"
	"database/sql"
	"time"
	"github.com/TechBowl-japan/go-stations/model"
	"github.com/mattn/go-sqlite3"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
    const (
        insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
        confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
    )

	now := time.Now()

    // トランザクションを開始
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    // INSERTクエリを実行
    res, err := tx.ExecContext(ctx, insert, subject, description, now , now)
    if err != nil {
        return nil, err // エラーをそのまま返す
    }

    // 挿入されたレコードのIDを取得
    id, err := res.LastInsertId()
    if err != nil {
        return nil, err
    }

    // 挿入したレコードを取得
    var todo model.TODO
    err = tx.QueryRowContext(ctx, confirm, id).Scan(
        &todo.Subject,
        &todo.Description,
        &todo.CreatedAt,
        &todo.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }

    // IDをセット
    todo.ID = id

    // トランザクションをコミット
    if err := tx.Commit(); err != nil {
        return nil, err
    }

    return &todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	return nil, nil
}

// UpdateTODO updates the TODO on DB.
// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ?, updated_at = ? WHERE id = ?`
		confirm = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	// ID が無効な場合、ErrNotFound を返す
	if id == 0 {
		return nil, &model.ErrNotFound{}
	}

	// Subject が空の場合、SQLite の制約エラーを模倣する
	if subject == "" {
		return nil, sqlite3.Error{Code: sqlite3.ErrConstraint}
	}

	// トランザクションを開始
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 現在時刻
	now := time.Now()

	// UPDATE クエリを実行
	res, err := tx.ExecContext(ctx, update, subject, description, now, id)
	if err != nil {
		return nil, err
	}

	// 更新された行数を確認
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		// 対象 ID のレコードが存在しない場合
		return nil, &model.ErrNotFound{}
	}

	// 更新後のレコードを取得
	var todo model.TODO
	err = tx.QueryRowContext(ctx, confirm, id).Scan(
		&todo.ID,
		&todo.Subject,
		&todo.Description,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// 対象 ID のレコードが存在しない場合
			return nil, &model.ErrNotFound{}
		}
		return nil, err
	}

	// トランザクションをコミット
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &todo, nil
}


// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}
