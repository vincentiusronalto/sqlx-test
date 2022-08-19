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

	err = tx.QueryRowContext(ctx, `INSERT INTO users (name) VALUES ($1) RETURNING id`, "user1").Scan(&lastInsertId)

	if err != nil {
		fmt.Println(err)
		tx.Rollback()
	}
	
	//lastInsertId :  24
	fmt.Println("lastInsertId : ", lastInsertId)

	result, err := tx.ExecContext(ctx, `INSERT INTO users (name) VALUES ($1)`, "user2")
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
	}
	

	id, errGetLastID := result.LastInsertId()
	if errGetLastID != nil {
		fmt.Println(err)
	}
	//lastInsertId2 :  0 -> cannot get since lastInserId not implemented in psql -> https://github.com/jmoiron/sqlx/issues/154
	fmt.Println("lastInsertId2 : ", id)
	tx.Commit()
}
