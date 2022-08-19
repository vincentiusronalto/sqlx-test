package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=alto password=alto dbname=test-migration sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	var lastInsertId int
	tx := db.MustBegin()

	// ================= 1) USING tx.QueryRowContext ================= SUCCESS
	err = tx.QueryRowContext(ctx, `INSERT INTO users (name) VALUES ($1) RETURNING id`, "user1").Scan(&lastInsertId)

	if err != nil {
		fmt.Println(err)
		tx.Rollback()
	}

	//lastInsertId :  29
	fmt.Println("lastInsertId : ", lastInsertId)

	// ================= 2) USING tx.GetContext ================= SUCCESS
	var lastInsertId2 int
	err = tx.GetContext(ctx, &lastInsertId2, `INSERT INTO users (name) VALUES ($1) RETURNING id`, "user2")

	if err != nil {
		fmt.Println(err)
		tx.Rollback()
	}

	//lastInsertId :  30
	fmt.Println("lastInsertId2 : ", lastInsertId2)

	// ================= 3) USING tx.ExecContext ======================= FAILED
	result, err := tx.ExecContext(ctx, `INSERT INTO users (name) VALUES ($1)`, "user3")
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
	}

	lastInsertId3 , errGetLastID := result.LastInsertId()
	if errGetLastID != nil {
		fmt.Println(err)
	}
	//lastInsertId2 :  0 -> cannot get since lastInserId not implemented in psql -> https://github.com/jmoiron/sqlx/issues/154
	fmt.Println("lastInsertId3 : ", lastInsertId3)
	tx.Commit()
}
