package main

import (
	"context"
	"database/sql" // SQLを投げる、エラーを返す、接続を管理するといった、共通の処理を扱う
	"fmt"
	"log"

	"github.com/pkg/errors"
	_ "modernc.org/sqlite" // アンダースコアを使用して、初期化(init)のみを実行
)

func run(ctx context.Context) error {
	db, err := sql.Open("sqlite", "example.db")
	if err != nil {
		return errors.Wrap(err, "failed to open database")
	}
	defer db.Close()

	_, err = db.ExecContext(ctx, `DROP TABLE IF EXISTS balances`)
	if err != nil {
		return errors.Wrap(err, "failed to create table")
	}

	// SQL命令には投げっぱなしのCRUD命令（Exec系）と、結果が見たい検索命令（Query系）がある
	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS balances (user_id INTEGER, balance INTEGER);`)
	if err != nil {
		return errors.Wrap(err, "failed to create table")
	}

	_, err = db.ExecContext(ctx, `INSERT INTO balances (user_id, balance) VALUES (41, 100);`)
	if err != nil {
		return errors.Wrap(err, "failed to insert initial data")
	}

	id := 41
	if err := updateBalance(ctx, db, id); err != nil {
		return errors.Wrap(err, "failed to process update sequence")
	}

	return nil
}

func updateBalance(ctx context.Context, db *sql.DB, id int) error {
	// dbとのコネクション（専属の1本）を確保
	conn, err := db.Conn(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get connection from pool")
	}
	defer conn.Close()

	// func (c *Conn) ExecContext(ctx context.Context, query string, args ...any) (Result, error)
	result, err := conn.ExecContext(ctx, `UPDATE balances SET balance = balance + 10 WHERE user_id = ?;`, id)
	if err != nil {
		return errors.Wrapf(err, "failed to execute update for user %d", id)
	}

	// sql.ResultはRowsAffected()メソッドを持つ
	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rows != 1 {
		return errors.Errorf("expected 1 row affected, got %d", rows)
	}

	fmt.Printf("成功! ユーザー%dの残高を更新しました。影響行数：%d\n", id, rows)
	return nil
}

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		// %+v を使うことで errors.Wrap で積み上げたエラーの詳細を出力
		log.Fatalf("Fatal error: %+v", err)
	}
}
