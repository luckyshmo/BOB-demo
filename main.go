package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aarondl/opt/omit"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
)

func NewDB() (*sqlx.DB, error) {
	db, err := sqlx.Open(
		"pgx",
		"postgres://postgres:example@127.0.0.1:6543/postgres?sslmode=disable",
	)
	if err != nil {
		return nil, fmt.Errorf("connect to API postgres: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping pg DB: %w", err)
	}

	return db, nil
}

type User struct {
	ID    int
	Name  string
	Email string
}

type UserSetter struct {
	ID    omit.Val[int]
	Name  omit.Val[string]
	Email omit.Val[string]
}

func main() {
	db, err := NewDB()
	if err != nil {
		log.Fatal(err)
	}

	bobDB := bob.New(db)

	var userTable = psql.NewTable[User, UserSetter]("public", "user_tb")
	userTable.Insert(context.TODO(), bobDB, UserSetter{
		ID:    omit.From(1),
		Name:  omit.From("Lol"),
		Email: omit.From("Kek"),
	})

	view := userTable.View

	q := view.Query(context.TODO(), bobDB, sm.Limit(10))
	users, err := q.All()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(users)
}
