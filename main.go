package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

func main() {
	ctx := context.Background()

	db, err := sql.Open("sqlite", "example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// SQL命令には投げっぱなしのCRUD命令、結果がみたい検索命令がある
	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS balances (user_id INTEGER, balance INTEGER);`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.ExecContext(ctx, `INSERT INTO balances (user_id, balance) VALUES (41, 100);`)
	if err != nil {
		log.Fatal(err)
	}
	// func (c *Conn) ExecContext(ctx context.Context, query string, args ...any) (Result, error)

	// dbとのコネクションを確保
	conn, err := db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	id := 41
	result, err := conn.ExecContext(ctx, `UPDATE balances SET balance = balance + 10 WHERE user_id = ?;`, id)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("成功! 更新された行数：%d\n", rows)

}
